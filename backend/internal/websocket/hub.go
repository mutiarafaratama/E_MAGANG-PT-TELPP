package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/telpp/emagang/internal/middleware"
	"github.com/telpp/emagang/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	HandshakeTimeout: 10 * time.Second,
}

// Client merepresentasikan satu koneksi WebSocket aktif
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID uuid.UUID
	Role   models.UserRole
}

// WsMessage adalah format standar semua pesan WebSocket
type WsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Hub mengelola semua koneksi aktif
type Hub struct {
	clients    map[uuid.UUID]*Client
	broadcast  chan *WsMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var GlobalHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		broadcast:  make(chan *WsMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("🔗 WS terhubung: user=%s role=%s (online: %d)", client.UserID, client.Role, h.OnlineCount())

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("🔌 WS terputus: user=%s (online: %d)", client.UserID, h.OnlineCount())

		case message := <-h.broadcast:
			data, _ := json.Marshal(message)
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// OnlineCount jumlah koneksi aktif saat ini
func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// IsOnline cek apakah user sedang terhubung
func (h *Hub) IsOnline(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// SendToUser kirim pesan ke user tertentu — return false jika tidak online
func (h *Hub) SendToUser(userID uuid.UUID, msgType string, payload interface{}) bool {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	msg := &WsMessage{Type: msgType, Payload: payload}
	data, err := json.Marshal(msg)
	if err != nil {
		return false
	}

	select {
	case client.Send <- data:
		return true
	default:
		h.mu.Lock()
		delete(h.clients, userID)
		close(client.Send)
		h.mu.Unlock()
		return false
	}
}

// SendToRole kirim pesan ke semua user dengan role tertentu
func (h *Hub) SendToRole(role models.UserRole, msgType string, payload interface{}) {
	msg := &WsMessage{Type: msgType, Payload: payload}
	data, _ := json.Marshal(msg)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.Role == role {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

// SendToRoles kirim ke beberapa role sekaligus
func (h *Hub) SendToRoles(roles []models.UserRole, msgType string, payload interface{}) {
	roleSet := make(map[models.UserRole]bool)
	for _, r := range roles {
		roleSet[r] = true
	}

	msg := &WsMessage{Type: msgType, Payload: payload}
	data, _ := json.Marshal(msg)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if roleSet[client.Role] {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// Ping untuk jaga koneksi tetap hidup
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		// Reset deadline setiap ada pesan masuk
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	}
}

// ServeWS — HTTP handler untuk upgrade ke WebSocket
// Token bisa dikirim via:
//   1. Header: Authorization: Bearer <token>
//   2. Query param: ?token=<token>  ← dipakai browser karena JS WS API tidak support custom header
func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token: coba dari query param dulu, lalu dari header
		token := c.Query("token")
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token diperlukan untuk koneksi WebSocket"})
			return
		}

		// Validasi token JWT
		claims, err := middleware.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid atau sudah kedaluwarsa"})
			return
		}

		// Upgrade ke WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := &Client{
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			UserID: claims.UserID,
			Role:   claims.Role,
		}

		hub.register <- client

		// Kirim konfirmasi koneksi berhasil
		welcome := &WsMessage{
			Type: "connected",
			Payload: map[string]interface{}{
				"user_id": claims.UserID,
				"role":    claims.Role,
				"message": "Koneksi WebSocket berhasil",
			},
		}
		if data, err := json.Marshal(welcome); err == nil {
			client.Send <- data
		}

		go client.writePump()
		go client.readPump()
	}
}
