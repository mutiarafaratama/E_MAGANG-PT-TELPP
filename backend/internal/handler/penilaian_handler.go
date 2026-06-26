package handler

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
)

type PenilaianHandler struct {
        repo *repository.PenilaianRepository
}

func NewPenilaianHandler(repo *repository.PenilaianRepository) *PenilaianHandler {
        return &PenilaianHandler{repo: repo}
}

// GET /api/penilaian/:pelaksanaan_id
func (h *PenilaianHandler) Get(c *gin.Context) {
        id, err := uuid.Parse(c.Param("pelaksanaan_id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }
        p, err := h.repo.GetByPelaksanaan(c.Request.Context(), id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        if p == nil {
                c.JSON(http.StatusOK, nil)
                return
        }
        c.JSON(http.StatusOK, p)
}

// POST /api/penilaian/:pelaksanaan_id
func (h *PenilaianHandler) Upsert(c *gin.Context) {
        pelaksanaanID, err := uuid.Parse(c.Param("pelaksanaan_id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_id", Message: "ID tidak valid"})
                return
        }

        var req models.PenilaianUpsertRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
                return
        }
        req.PelaksanaanID = pelaksanaanID

        dinilaiOleh := middleware.GetUserID(c)

        result, err := h.repo.Upsert(c.Request.Context(), req, dinilaiOleh)
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        c.JSON(http.StatusOK, models.SuccessResponse{Message: "Penilaian berhasil disimpan", Data: result})
}

// GET /api/penilaian
func (h *PenilaianHandler) List(c *gin.Context) {
        list, err := h.repo.ListPerluDinilai(c.Request.Context())
        if err != nil {
                c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "server_error", Message: err.Error()})
                return
        }
        if list == nil {
                list = []models.PenilaianListItem{}
        }
        c.JSON(http.StatusOK, list)
}

