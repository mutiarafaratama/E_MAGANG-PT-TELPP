package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

// GET /api/admin/stats — statistik dashboard
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.userRepo.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: stats})
}

// GET /api/admin/users — daftar semua user
func (h *AdminHandler) GetUsers(c *gin.Context) {
	role := c.Query("role")
	page := queryInt(c, "page", 1)
	limit := queryInt(c, "limit", 20)

	users, total, err := h.userRepo.FindAll(c.Request.Context(), role, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data: users, Total: total, Page: page, Limit: limit,
	})
}

// POST /api/admin/users — buat akun HRD baru
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req struct {
		NamaLengkap string          `json:"nama_lengkap" binding:"required"`
		Email       string          `json:"email" binding:"required,email"`
		Password    string          `json:"password" binding:"required,min=8"`
		Role        models.UserRole `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}

	if req.Role != models.RoleHRD && req.Role != models.RoleAdmin {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_role", Message: "Role hanya boleh hrd atau admin"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	user := &models.User{
		NamaLengkap:  req.NamaLengkap,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
	}

	if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "create_failed", Message: "Email sudah terdaftar"})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Akun berhasil dibuat",
		Data: models.UserPublic{
			ID:          user.ID,
			NamaLengkap: user.NamaLengkap,
			Email:       user.Email,
			Role:        user.Role,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
		},
	})
}

// PATCH /api/admin/users/:id/toggle — aktifkan/nonaktifkan user
func (h *AdminHandler) ToggleUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	c.ShouldBindJSON(&req)

	if err := h.userRepo.UpdateActive(c.Request.Context(), id, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}

	status := "dinonaktifkan"
	if req.IsActive {
		status = "diaktifkan"
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Akun berhasil " + status})
}

// GET /api/admin/hrd-list — daftar HRD untuk assign
func (h *AdminHandler) GetHRDList(c *gin.Context) {
	list, err := h.userRepo.FindHRDList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

func parseTimeStr(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
