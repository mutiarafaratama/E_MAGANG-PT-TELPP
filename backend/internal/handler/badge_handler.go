package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/telpp/emagang/internal/database"
	"github.com/telpp/emagang/internal/middleware"
	"github.com/telpp/emagang/internal/models"
)

type BadgeHandler struct{}

func NewBadgeHandler() *BadgeHandler { return &BadgeHandler{} }

// GET /api/badge/hrd — badge counts untuk semua menu HRD
func (h *BadgeHandler) GetHRD(c *gin.Context) {
	ctx := c.Request.Context()
	db := database.DB

	var verifikasi, penempatan, izinPending, laporanPending, penilaianSiap, chatUnread, perpanjanganPending int

	db.QueryRow(ctx, `SELECT COUNT(*) FROM pengajuan_magang WHERE status IN ('diajukan','menunggu_verifikasi')`).
		Scan(&verifikasi)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM pengajuan_magang p
		WHERE p.status = 'diterima'
		AND NOT EXISTS (SELECT 1 FROM pelaksanaan_magang pm WHERE pm.pengajuan_id = p.id)`).
		Scan(&penempatan)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM izin_sakit_request WHERE status = 'pending'`).
		Scan(&izinPending)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM laporan_magang WHERE status = 'menunggu_review'`).
		Scan(&laporanPending)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM pelaksanaan_magang WHERE status = 'penilaian'`).
		Scan(&penilaianSiap)

	db.QueryRow(ctx, `SELECT COUNT(DISTINCT t.id) FROM chat_tiket t
		JOIN chat_pesan p ON p.tiket_id = t.id
		WHERE p.sender_id = t.user_id AND p.is_read = false`).
		Scan(&chatUnread)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM perpanjangan_magang WHERE status = 'menunggu'`).
		Scan(&perpanjanganPending)

	c.JSON(http.StatusOK, models.SuccessResponse{Data: map[string]any{
		"verifikasi_berkas":    verifikasi,
		"menunggu_penempatan":  penempatan,
		"izin_sakit_pending":   izinPending,
		"laporan_pending":      laporanPending,
		"penilaian_siap":       penilaianSiap,
		"chat_unread":          chatUnread,
		"perpanjangan_pending": perpanjanganPending,
	}})
}

// GET /api/badge/peserta — badge counts untuk menu peserta yang sedang login
func (h *BadgeHandler) GetPeserta(c *gin.Context) {
	ctx := c.Request.Context()
	db := database.DB
	userID := middleware.GetUserID(c)

	var laporanRevisi, nilaiBaru, sertifikatBaru, chatUnread int

	db.QueryRow(ctx, `SELECT COUNT(*) FROM laporan_magang lm
		JOIN pelaksanaan_magang pm ON pm.id = lm.pelaksanaan_id
		WHERE pm.user_id = $1
		AND lm.status = 'revisi'
		AND lm.versi = (SELECT MAX(versi) FROM laporan_magang WHERE pelaksanaan_id = pm.id)`,
		userID).Scan(&laporanRevisi)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM pelaksanaan_magang
		WHERE user_id = $1 AND nilai IS NOT NULL`, userID).Scan(&nilaiBaru)

	db.QueryRow(ctx, `SELECT COUNT(*) FROM pelaksanaan_magang
		WHERE user_id = $1 AND sertifikat_generated = true`, userID).Scan(&sertifikatBaru)

	db.QueryRow(ctx, `SELECT COUNT(DISTINCT t.id) FROM chat_tiket t
		JOIN chat_pesan p ON p.tiket_id = t.id
		WHERE t.user_id = $1 AND p.sender_id != $1 AND p.is_read = false`, userID).Scan(&chatUnread)

	c.JSON(http.StatusOK, models.SuccessResponse{Data: map[string]any{
		"laporan_revisi":  laporanRevisi,
		"nilai_baru":      nilaiBaru,
		"sertifikat_baru": sertifikatBaru,
		"chat_unread":     chatUnread,
	}})
}
