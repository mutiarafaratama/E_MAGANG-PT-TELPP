package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telpp/emagang/internal/models"
	"github.com/telpp/emagang/internal/repository"
	"github.com/telpp/emagang/internal/service"
)

type KnowledgeHandler struct {
	repo   *repository.KnowledgeRepository
	nlpSvc *service.NLPService
}

func NewKnowledgeHandler(repo *repository.KnowledgeRepository, nlpSvc *service.NLPService) *KnowledgeHandler {
	return &KnowledgeHandler{repo: repo, nlpSvc: nlpSvc}
}

// GetAll — admin: list semua knowledge (aktif dan tidak)
func (h *KnowledgeHandler) GetAll(c *gin.Context) {
	list, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data knowledge"})
		return
	}
	if list == nil {
		list = []models.ChatKnowledge{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Create — admin: buat knowledge baru
func (h *KnowledgeHandler) Create(c *gin.Context) {
	var req models.ChatKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	k, err := h.repo.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat knowledge"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": k})
}

// Update — admin: update knowledge
func (h *KnowledgeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	var req models.ChatKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	k, err := h.repo.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal update knowledge"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": k})
}

// Delete — admin: hapus knowledge
func (h *KnowledgeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal hapus knowledge"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "knowledge berhasil dihapus"})
}

// Toggle — admin: aktifkan/nonaktifkan knowledge
func (h *KnowledgeHandler) Toggle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}
	k, err := h.repo.ToggleActive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal toggle knowledge"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": k})
}

// QueryNLP — public & logged-in: tanya chatbot
func (h *KnowledgeHandler) QueryNLP(c *gin.Context) {
	var req models.NLPQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	result := h.nlpSvc.Query(c.Request.Context(), req.Pesan)
	c.JSON(http.StatusOK, gin.H{"data": result})
}
