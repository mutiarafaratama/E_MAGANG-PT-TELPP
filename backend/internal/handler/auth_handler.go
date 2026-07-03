package handler

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/service"
)

type AuthHandler struct {
        svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
        return &AuthHandler{svc: svc}
}

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
        var req models.RegisterRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        resp, err := h.svc.Register(c.Request.Context(), req)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "register_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Registrasi berhasil",
                Data:    resp,
        })
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
        var req models.LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        resp, err := h.svc.Login(c.Request.Context(), req)
        if err != nil {
                c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "login_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{
                Message: "Login berhasil",
                Data:    resp,
        })
}

// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
        userID := middleware.GetUserID(c)
        profile, err := h.svc.GetProfile(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: profile})
}

// POST /api/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
        var req models.ForgotPasswordRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        if err := h.svc.ForgotPassword(c.Request.Context(), req.Email); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "not_peserta", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{
                Message: "Tautan reset kata sandi telah dikirimkan ke email Anda",
        })
}

// POST /api/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
        var req models.ResetPasswordRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        if err := h.svc.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "reset_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kata sandi berhasil diperbarui"})
}

// POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
        var req struct {
                OldPassword string `json:"old_password" binding:"required"`
                NewPassword string `json:"new_password" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        userID := middleware.GetUserID(c)
        if err := h.svc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "change_password_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Password berhasil diubah"})
}
