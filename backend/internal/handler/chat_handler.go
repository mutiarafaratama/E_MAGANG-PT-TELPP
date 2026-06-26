package handler

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/service"
)

type ChatHandler struct {
        svc *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
        return &ChatHandler{svc: svc}
}

// POST /api/chat/tiket — peserta buat tiket baru
func (h *ChatHandler) BuatTiket(c *gin.Context) {
        var req models.ChatBuatTiketRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        userID := middleware.GetUserID(c)
        tiket, err := h.svc.BuatTiket(c.Request.Context(), userID, req)
        if err != nil {
                status := http.StatusInternalServerError
                if err.Error() != "" {
                        status = http.StatusConflict
                }
                c.JSON(status, models.ErrorResponse{Error: "buat_tiket_gagal", Message: err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Tiket berhasil dibuat",
                Data:    tiket,
        })
}

// GET /api/chat/tiket/saya — peserta lihat tiket sendiri
func (h *ChatHandler) GetTiketSaya(c *gin.Context) {
        userID := middleware.GetUserID(c)
        list, err := h.svc.GetTiketSaya(c.Request.Context(), userID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Data: list})
}

// GET /api/chat/tiket — HRD/Admin lihat semua tiket
func (h *ChatHandler) GetAllTiket(c *gin.Context) {
        status   := c.Query("status")
        kategori := c.Query("kategori")
        page     := queryInt(c, "page", 1)
        limit    := queryInt(c, "limit", 20)

        list, total, err := h.svc.GetAllTiket(c.Request.Context(), status, kategori, page, limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.PaginatedResponse{
                Data: list, Total: total, Page: page, Limit: limit,
        })
}

// GET /api/chat/tiket/:id/pesan — ambil pesan dalam tiket
func (h *ChatHandler) GetPesan(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        userID := middleware.GetUserID(c)
        role   := middleware.GetUserRole(c)

        pesan, err := h.svc.GetPesan(c.Request.Context(), id, userID, role)
        if err != nil {
                c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "forbidden", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Data: pesan})
}

// POST /api/chat/tiket/:id/pesan — kirim pesan
func (h *ChatHandler) KirimPesan(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req models.ChatKirimPesanRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        userID := middleware.GetUserID(c)
        role   := middleware.GetUserRole(c)

        pesan, err := h.svc.KirimPesan(c.Request.Context(), id, userID, role, req)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "send_failed", Message: err.Error()})
                return
        }

        c.JSON(http.StatusCreated, models.SuccessResponse{
                Message: "Pesan terkirim",
                Data:    pesan,
        })
}

// PATCH /api/chat/tiket/:id/ambil — HRD ambil/assign tiket ke diri sendiri
func (h *ChatHandler) AmbilTiket(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        hrdID := middleware.GetUserID(c)
        tiket, err := h.svc.AmbilTiket(c.Request.Context(), id, hrdID)
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "ambil_gagal", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{
                Message: "Tiket berhasil diambil",
                Data:    tiket,
        })
}

// PATCH /api/chat/tiket/:id/status — HRD update status tiket
func (h *ChatHandler) UpdateStatus(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req struct {
                Status     models.StatusTiket `json:"status" binding:"required"`
                AssignedTo *string            `json:"assigned_to"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }

        var assignedTo *uuid.UUID
        if req.AssignedTo != nil {
                pid, err := uuid.Parse(*req.AssignedTo)
                if err == nil {
                        assignedTo = &pid
                }
        }

        if err := h.svc.UpdateStatusTiket(c.Request.Context(), id, req.Status, assignedTo); err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }

        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Status tiket diperbarui"})
}

