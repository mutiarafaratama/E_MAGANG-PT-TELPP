package handler

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        "github.com/telpp/emagang/internal/service"
)

type NotifikasiHandler struct {
        svc      *service.NotifikasiService
        chatRepo *repository.ChatRepository
}

func NewNotifikasiHandler(svc *service.NotifikasiService, chatRepo *repository.ChatRepository) *NotifikasiHandler {
        return &NotifikasiHandler{svc: svc, chatRepo: chatRepo}
}

// GET /api/notifikasi — list notifikasi user
func (h *NotifikasiHandler) GetAll(c *gin.Context) {
        userID := middleware.GetUserID(c)
        page := queryInt(c, "page", 1)
        limit := queryInt(c, "limit", 20)

        list, total, err := h.svc.GetAll(c.Request.Context(), userID, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.PaginatedResponse{
                Data: list, Total: total, Page: page, Limit: limit,
        })
}

// PATCH /api/notifikasi/:id/read
func (h *NotifikasiHandler) MarkRead(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        userID := middleware.GetUserID(c)
        h.svc.MarkRead(c.Request.Context(), id, userID)
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Notifikasi ditandai sudah dibaca"})
}

// PATCH /api/notifikasi/read-all
func (h *NotifikasiHandler) MarkAllRead(c *gin.Context) {
        userID := middleware.GetUserID(c)
        h.svc.MarkAllRead(c.Request.Context(), userID)
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Semua notifikasi ditandai sudah dibaca"})
}

// GET /api/notifikasi/badge — jumlah unread + chat menunggu (untuk sidebar badge)
func (h *NotifikasiHandler) GetBadge(c *gin.Context) {
        userID := middleware.GetUserID(c)
        badge := h.svc.GetBadgeCount(c.Request.Context(), userID, h.chatRepo)
        c.JSON(http.StatusOK, models.SuccessResponse{Data: badge})
}

// POST /api/notifikasi/fcm-token — simpan FCM token perangkat
func (h *NotifikasiHandler) SaveFCMToken(c *gin.Context) {
        var req struct {
                Token string `json:"token" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }
        userID := middleware.GetUserID(c)
        if err := h.svc.SaveFCMToken(c.Request.Context(), userID, req.Token); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: "Gagal menyimpan FCM token"})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "FCM token berhasil disimpan"})
}

// Helper untuk parse query int
func queryInt(c *gin.Context, key string, def int) int {
        val := c.Query(key)
        if val == "" {
                return def
        }
        n, err := strconv.Atoi(val)
        if err != nil {
                return def
        }
        return n
}
