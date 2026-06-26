package service

import (
        "bytes"
        "context"
        "crypto/rsa"
        "crypto/x509"
        "encoding/json"
        "encoding/pem"
        "fmt"
        "io"
        "log"
        "net/http"
        "net/url"
        "strings"
        "sync"
        "time"

        "github.com/golang-jwt/jwt/v5"
)

const (
        fcmV1Endpoint  = "https://fcm.googleapis.com/v1/projects/%s/messages:send"
        googleTokenURL = "https://oauth2.googleapis.com/token"
        fcmScope       = "https://www.googleapis.com/auth/firebase.messaging"
)

type serviceAccountJSON struct {
        Type        string `json:"type"`
        ProjectID   string `json:"project_id"`
        PrivateKey  string `json:"private_key"`
        ClientEmail string `json:"client_email"`
        TokenURI    string `json:"token_uri"`
}

// FCMService kirim push notification ke perangkat via Firebase Cloud Messaging V1 API (ServiceAccount).
// Jika kredensial tidak dikonfigurasi semua operasi adalah no-op (graceful degradation).
type FCMService struct {
        projectID  string
        sa         *serviceAccountJSON
        privateKey *rsa.PrivateKey
        client     *http.Client
        mu         sync.Mutex
        token      string
        tokenExp   time.Time
}

func NewFCMService(projectID, serviceAccountJSONStr string) *FCMService {
        svc := &FCMService{
                projectID: projectID,
                client:    &http.Client{Timeout: 10 * time.Second},
        }
        if serviceAccountJSONStr == "" || projectID == "" {
                return svc
        }
        var sa serviceAccountJSON
        if err := json.Unmarshal([]byte(serviceAccountJSONStr), &sa); err != nil {
                log.Printf("[FCM] gagal parse service account JSON: %v", err)
                return svc
        }
        key, err := parseRSAPrivateKey(sa.PrivateKey)
        if err != nil {
                log.Printf("[FCM] gagal parse private key: %v", err)
                return svc
        }
        svc.sa = &sa
        svc.privateKey = key
        log.Printf("[FCM] V1 service diinisialisasi untuk project: %s", projectID)
        return svc
}

func parseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
        // Env var/secret sering menyimpan \n sebagai literal backslash-n — normalkan dulu
        pemStr = strings.ReplaceAll(pemStr, `\n`, "\n")
        block, _ := pem.Decode([]byte(pemStr))
        if block == nil {
                return nil, fmt.Errorf("invalid PEM block")
        }
        key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
        if err == nil {
                rsaKey, ok := key.(*rsa.PrivateKey)
                if !ok {
                        return nil, fmt.Errorf("bukan RSA key")
                }
                return rsaKey, nil
        }
        return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func (f *FCMService) Enabled() bool {
        return f.projectID != "" && f.sa != nil && f.privateKey != nil
}

func (f *FCMService) getAccessToken(ctx context.Context) (string, error) {
        f.mu.Lock()
        defer f.mu.Unlock()

        if f.token != "" && time.Now().Before(f.tokenExp.Add(-60*time.Second)) {
                return f.token, nil
        }

        now := time.Now()
        tokenURI := f.sa.TokenURI
        if tokenURI == "" {
                tokenURI = googleTokenURL
        }

        claims := jwt.MapClaims{
                "iss":   f.sa.ClientEmail,
                "sub":   f.sa.ClientEmail,
                "aud":   tokenURI,
                "scope": fcmScope,
                "iat":   now.Unix(),
                "exp":   now.Add(time.Hour).Unix(),
        }
        signedJWT, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(f.privateKey)
        if err != nil {
                return "", fmt.Errorf("sign JWT: %w", err)
        }

        resp, err := f.client.PostForm(tokenURI, url.Values{
                "grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
                "assertion":  {signedJWT},
        })
        if err != nil {
                return "", fmt.Errorf("token request: %w", err)
        }
        defer resp.Body.Close()

        body, _ := io.ReadAll(resp.Body)
        var result struct {
                AccessToken string `json:"access_token"`
                ExpiresIn   int    `json:"expires_in"`
        }
        if err := json.Unmarshal(body, &result); err != nil || result.AccessToken == "" {
                return "", fmt.Errorf("token response invalid: %s", string(body))
        }

        f.token = result.AccessToken
        f.tokenExp = now.Add(time.Duration(result.ExpiresIn) * time.Second)
        return f.token, nil
}

type fcmV1Request struct {
        Message fcmV1Message `json:"message"`
}

type fcmV1Message struct {
        Token        string            `json:"token"`
        Notification fcmNotification   `json:"notification"`
        Data         map[string]string `json:"data,omitempty"`
        WebPush      *fcmWebPush       `json:"webpush,omitempty"`
}

type fcmNotification struct {
        Title string `json:"title"`
        Body  string `json:"body"`
}

type fcmWebPush struct {
        Notification *fcmWebPushNotif `json:"notification,omitempty"`
}

type fcmWebPushNotif struct {
        Title string `json:"title"`
        Body  string `json:"body"`
        Icon  string `json:"icon,omitempty"`
}

// Send kirim push notification ke satu FCM token via V1 API.
func (f *FCMService) Send(ctx context.Context, fcmToken, judul, pesan string, data map[string]string) {
        if !f.Enabled() || fcmToken == "" {
                return
        }

        accessToken, err := f.getAccessToken(ctx)
        if err != nil {
                log.Printf("[FCM V1] gagal get access token: %v", err)
                return
        }

        payload := fcmV1Request{
                Message: fcmV1Message{
                        Token:        fcmToken,
                        Notification: fcmNotification{Title: judul, Body: pesan},
                        Data:         data,
                        WebPush: &fcmWebPush{
                                Notification: &fcmWebPushNotif{
                                        Title: judul,
                                        Body:  pesan,
                                        Icon:  "/logo.png",
                                },
                        },
                },
        }

        body, err := json.Marshal(payload)
        if err != nil {
                log.Printf("[FCM V1] marshal error: %v", err)
                return
        }

        endpoint := fmt.Sprintf(fcmV1Endpoint, f.projectID)
        req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
        if err != nil {
                log.Printf("[FCM V1] request error: %v", err)
                return
        }
        req.Header.Set("Authorization", "Bearer "+accessToken)
        req.Header.Set("Content-Type", "application/json")

        resp, err := f.client.Do(req)
        if err != nil {
                log.Printf("[FCM V1] send error: %v", err)
                return
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 400 {
                respBody, _ := io.ReadAll(resp.Body)
                preview := fcmToken
                if len(preview) > 8 {
                        preview = preview[:8]
                }
                log.Printf("[FCM V1] server returned %d untuk token %s...: %s", resp.StatusCode, preview, string(respBody))
        }
}
