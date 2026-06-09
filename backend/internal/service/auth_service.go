package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/telpp/emagang/internal/middleware"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	passwordRegexLetter  = regexp.MustCompile(`[A-Za-z]`)
	passwordRegexDigit   = regexp.MustCompile(`[0-9]`)
	passwordRegexSpecial = regexp.MustCompile(`[!@#$%^&*()_+=\.,><?/]`)
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func isPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}
	return passwordRegexLetter.MatchString(password) &&
		passwordRegexDigit.MatchString(password) &&
		passwordRegexSpecial.MatchString(password)
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Validasi password
	if !isPasswordValid(req.Password) {
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
			ID:          user.ID,
			NamaLengkap: user.NamaLengkap,
			Email:       user.Email,
			Role:        user.Role,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserPublic, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &models.UserPublic{
		ID:          user.ID,
		NamaLengkap: user.NamaLengkap,
		Email:       user.Email,
		Role:        user.Role,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
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

	if !isPasswordValid(newPass) {
		return errors.New("password baru harus mengandung huruf, angka, dan karakter spesial (!@#$%^&*()_+=.,><?/)")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(newPass), 12)
	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}
