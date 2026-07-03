package service

import (
        "context"
        "crypto/sha256"
        "errors"
        "fmt"
        "time"
        "unicode"

        "github.com/google/uuid"
        "github.com/jackc/pgx/v5"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        "golang.org/x/crypto/bcrypt"
)

// validatePassword memeriksa bahwa password mengandung huruf, angka, dan karakter spesial
func validatePassword(s string) bool {
        if len(s) < 8 {
                return false
        }
        var hasLetter, hasDigit, hasSpecial bool
        special := "!@#$%^&*()_+=.,><?/"
        for _, c := range s {
                switch {
                case unicode.IsLetter(c):
                        hasLetter = true
                case unicode.IsDigit(c):
                        hasDigit = true
                default:
                        for _, sc := range special {
                                if c == sc {
                                        hasSpecial = true
                                        break
                                }
                        }
                }
        }
        return hasLetter && hasDigit && hasSpecial
}

type AuthService struct {
        userRepo *repository.UserRepository
        emailSvc *EmailService
}

func NewAuthService(userRepo *repository.UserRepository, emailSvc *EmailService) *AuthService {
        return &AuthService{userRepo: userRepo, emailSvc: emailSvc}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
        // Validasi password
        if !validatePassword(req.Password) {
                return nil, errors.New("password harus mengandung huruf, angka, dan karakter spesial (!@#$%^&*()_+=.,><?/)")
        }

        // Cek email sudah terdaftar
        existing, err := s.userRepo.FindByEmail(ctx, req.Email)
        if err == nil && existing != nil {
                return nil, errors.New("email sudah terdaftar")
        }

        // Hash password
        hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
        if err != nil {
                return nil, fmt.Errorf("gagal hash password: %w", err)
        }

        user := &models.User{
                NamaLengkap:  req.NamaLengkap,
                Email:        req.Email,
                PasswordHash: string(hash),
                Role:         models.RolePeserta,
        }

        if err := s.userRepo.Create(ctx, user); err != nil {
                return nil, fmt.Errorf("gagal membuat akun: %w", err)
        }

        return s.generateTokenPair(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
        user, err := s.userRepo.FindByEmail(ctx, req.Email)
        if err != nil {
                if errors.Is(err, pgx.ErrNoRows) {
                        return nil, errors.New("email atau password salah")
                }
                return nil, err
        }

        if !user.IsActive {
                return nil, errors.New("akun Anda telah dinonaktifkan. Hubungi admin")
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
                return nil, errors.New("email atau password salah")
        }

        return s.generateTokenPair(ctx, user)
}

func (s *AuthService) generateTokenPair(ctx context.Context, user *models.User) (*models.AuthResponse, error) {
        accessToken, err := middleware.GenerateAccessToken(user.ID, user.Email, user.Role)
        if err != nil {
                return nil, fmt.Errorf("gagal generate token: %w", err)
        }

        refreshToken, err := middleware.GenerateRefreshToken(user.ID)
        if err != nil {
                return nil, fmt.Errorf("gagal generate refresh token: %w", err)
        }

        // Simpan hash refresh token
        hash := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))
        s.userRepo.SaveRefreshToken(ctx, user.ID, hash)

        return &models.AuthResponse{
                AccessToken:  accessToken,
                RefreshToken: refreshToken,
                User: models.UserPublic{
                        ID:              user.ID,
                        NamaLengkap:     user.NamaLengkap,
                        Email:           user.Email,
                        Role:            user.Role,
                        IsActive:        user.IsActive,
                        PasswordChanged: user.PasswordChanged,
                        CreatedAt:       user.CreatedAt,
                },
        }, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserPublic, error) {
        user, err := s.userRepo.FindByID(ctx, userID)
        if err != nil {
                return nil, err
        }
        return &models.UserPublic{
                ID:              user.ID,
                NamaLengkap:     user.NamaLengkap,
                Email:           user.Email,
                Role:            user.Role,
                IsActive:        user.IsActive,
                PasswordChanged: user.PasswordChanged,
                CreatedAt:       user.CreatedAt,
        }, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPass, newPass string) error {
        user, err := s.userRepo.FindByID(ctx, userID)
        if err != nil {
                return err
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPass)); err != nil {
                return errors.New("password lama tidak sesuai")
        }

        if !validatePassword(newPass) {
                return errors.New("password baru harus mengandung huruf, angka, dan karakter spesial (!@#$%^&*()_+=.,><?/)")
        }

        hash, _ := bcrypt.GenerateFromPassword([]byte(newPass), 12)
        return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

// ForgotPassword — generate token dan kirim email reset kata sandi (hanya peserta magang)
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
        user, err := s.userRepo.FindByEmail(ctx, email)
        if err != nil {
                return errors.New("Email tidak terdaftar sebagai peserta magang pada sistem e-Magang TELPP")
        }

        if user.Role != models.RolePeserta {
                return errors.New("Anda belum menjadi peserta magang pada sistem e-Magang TELPP")
        }

        // Generate token acak (UUID) + hash untuk disimpan di DB
        rawToken := uuid.New().String()
        tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawToken)))
        expiredAt := time.Now().Add(1 * time.Hour)

        if err := s.userRepo.SavePasswordResetToken(ctx, user.ID, tokenHash, expiredAt); err != nil {
                return fmt.Errorf("gagal menyimpan token: %w", err)
        }

        resetURL := frontendURL() + "/reset-kata-sandi?token=" + rawToken
        _ = s.emailSvc.KirimResetPassword(user.Email, user.NamaLengkap, resetURL)

        return nil
}

// ResetPassword — validasi token dan update password user
func (s *AuthService) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
        if !validatePassword(newPassword) {
                return errors.New("password harus mengandung huruf, angka, dan karakter spesial (!@#$%^&*()_+=.,><?/)")
        }

        tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawToken)))
        t, err := s.userRepo.FindPasswordResetToken(ctx, tokenHash)
        if err != nil {
                return errors.New("tautan reset tidak valid atau sudah kedaluwarsa")
        }

        if t.Used {
                return errors.New("tautan reset sudah pernah digunakan")
        }

        if time.Now().After(t.ExpiredAt) {
                return errors.New("tautan reset sudah kedaluwarsa. Minta tautan baru")
        }

        hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
        if err != nil {
                return fmt.Errorf("gagal hash password: %w", err)
        }

        if err := s.userRepo.UpdatePassword(ctx, t.UserID, string(hash)); err != nil {
                return fmt.Errorf("gagal memperbarui password: %w", err)
        }

        _ = s.userRepo.MarkResetTokenUsed(ctx, tokenHash)
        return nil
}
