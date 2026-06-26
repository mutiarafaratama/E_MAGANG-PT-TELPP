package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/handler"
)

// RegisterKnowledgeRoutes mendaftarkan route knowledge base chatbot.
//
//	public : /api  (tanpa auth — untuk chatbot publik)
//	admin  : /api/admin  (hanya admin)
func RegisterKnowledgeRoutes(public, admin *gin.RouterGroup, h *handler.KnowledgeHandler) {
	// Public — query NLP (anonymous maupun logged-in)
	public.POST("/chatbot/query", h.QueryNLP)

	// Admin — kelola knowledge base
	admin.GET("/knowledge", h.GetAll)
	admin.POST("/knowledge", h.Create)
	admin.PUT("/knowledge/:id", h.Update)
	admin.DELETE("/knowledge/:id", h.Delete)
	admin.PATCH("/knowledge/:id/toggle", h.Toggle)
}
