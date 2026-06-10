package router

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "github.com/telpp/emagang/internal/database"
        "github.com/telpp/emagang/internal/handler"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        "github.com/telpp/emagang/internal/service"
        ws "github.com/telpp/emagang/internal/websocket"
)

func Setup() *gin.Engine {
        r := gin.Default()

        // ============================================================
        // MIDDLEWARE GLOBAL
        // ============================================================
        r.Use(middleware.CORS())
        r.MaxMultipartMemory = 110 << 20 // 110MB

        // ============================================================
        // INISIALISASI SEMUA LAYER
        // ============================================================

        // Repositories
        userRepo        := repository.NewUserRepository(database.DB)
        pengajuanRepo   := repository.NewPengajuanRepository(database.DB)
        dokumenRepo     := repository.NewDokumenRepository(database.DB)
        pelaksanaanRepo := repository.NewPelaksanaanRepository(database.DB)
        absensiRepo     := repository.NewAbsensiRepository(database.DB)
        chatRepo        := repository.NewChatRepository(database.DB)
        notifRepo       := repository.NewNotifikasiRepository(database.DB)
        landingRepo     := repository.NewLandingRepository(database.DB)

        // Set global DB pointer
        repository.SetDB(database.DB)

        // Services
        hub           := ws.GlobalHub
        notifSvc      := service.NewNotifikasiService(notifRepo, hub)
        authSvc       := service.NewAuthService(userRepo)
        pengajuanSvc  := service.NewPengajuanService(pengajuanRepo, notifSvc, userRepo)
        chatSvc       := service.NewChatService(chatRepo, notifSvc, userRepo, hub)
        sertifikatSvc := service.NewSertifikatService(pelaksanaanRepo, pengajuanRepo, notifSvc)

        // Handlers
        authH        := handler.NewAuthHandler(authSvc)
        pengajuanH   := handler.NewPengajuanHandler(pengajuanSvc)
        dokumenH     := handler.NewDokumenHandler(dokumenRepo)
        pelaksanaanH := handler.NewPelaksanaanHandler(pelaksanaanRepo, pengajuanRepo, sertifikatSvc)
        absensiH     := handler.NewAbsensiHandler(absensiRepo, pelaksanaanRepo, pengajuanRepo)
        chatH        := handler.NewChatHandler(chatSvc)
        notifH       := handler.NewNotifikasiHandler(notifSvc, chatRepo)
        landingH     := handler.NewLandingHandler(landingRepo)
        adminH       := handler.NewAdminHandler(userRepo)

        // ============================================================
        // HEALTH CHECK
        // ============================================================
        r.GET("/api/health", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"status": "ok", "app": "e-Magang TELPP"})
        })

        // DEBUG: Cek dokumen di database (hapus setelah production)
        r.GET("/api/debug/dokumen", func(c *gin.Context) {
                rows, err := database.DB.Query(c.Request.Context(), "SELECT id, pengajuan_id, user_id, jenis, nama_file FROM dokumen LIMIT 20")
                if err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                        return
                }
                defer rows.Close()

                var docs []gin.H
                for rows.Next() {
                        var id, pengajuanID, userID uuid.UUID
                        var jenis, namaFile string
                        var pengajuanIDPtr *uuid.UUID

                        if err := rows.Scan(&id, &pengajuanIDPtr, &userID, &jenis, &namaFile); err != nil {
                                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                                return
                        }

                        docs = append(docs, gin.H{
                                "id":           id,
                                "pengajuan_id": pengajuanIDPtr,
                                "user_id":      userID,
                                "jenis":        jenis,
                                "nama_file":    namaFile,
                        })
                }

                c.JSON(http.StatusOK, gin.H{"data": docs, "count": len(docs)})
        })

        // ============================================================
        // PUBLIC ROUTES (tanpa login)
        // ============================================================
        public := r.Group("/api")
        {
                // Auth
                public.POST("/auth/register", authH.Register)
                public.POST("/auth/login", authH.Login)

                // Landing page (publik)
                public.GET("/landing/content", landingH.GetContent)
                public.GET("/landing/faq", landingH.GetFAQ)
                public.GET("/landing/periode", landingH.GetPeriodeAktif)

                // WebSocket — autentikasi ditangani sendiri via query param ?token=
                // Browser JS tidak bisa kirim custom header saat upgrade WS
                // Contoh: new WebSocket("ws://localhost:8080/api/ws?token=<jwt>")
                public.GET("/ws", ws.ServeWS(hub))
        }

        // ============================================================
        // PROTECTED ROUTES (wajib login)
        // ============================================================
        api := r.Group("/api")
        api.Use(middleware.AuthRequired())
        {
                // ============================================================
                // SETUP ROUTES MODULAR
                // ============================================================
                SetupPengajuanRoutes(api, pengajuanH)
                SetupDokumenRoutes(api, dokumenH)

                // Auth
                api.GET("/auth/me", authH.Me)
                api.POST("/auth/change-password", authH.ChangePassword)

                // Notifikasi
                api.GET("/notifikasi", notifH.GetAll)
                api.PATCH("/notifikasi/:id/read", notifH.MarkRead)
                api.PATCH("/notifikasi/read-all", notifH.MarkAllRead)
                api.GET("/notifikasi/badge", notifH.GetBadge)

                // --------------------------------------------------------
                // PESERTA ROUTES
                // --------------------------------------------------------
                peserta := api.Group("/peserta")
                peserta.Use(middleware.RoleRequired(models.RolePeserta))
                {
                        // Pelaksanaan
                        peserta.GET("/pelaksanaan/saya", pelaksanaanH.GetMy)
                        peserta.GET("/pelaksanaan/:id/sertifikat/download", pelaksanaanH.DownloadSertifikat)

                        // Absensi
                        peserta.POST("/absensi/checkin", absensiH.CheckIn)
                        peserta.PATCH("/absensi/checkout", absensiH.CheckOut)
                        peserta.GET("/absensi/saya", absensiH.GetMy)
                        peserta.GET("/absensi/saya/pdf", absensiH.DownloadPDF)

                        // Chat (peserta buat tiket)
                        peserta.POST("/chat/tiket", chatH.BuatTiket)
                        peserta.GET("/chat/tiket/saya", chatH.GetTiketSaya)
                }

                // --------------------------------------------------------
                // SHARED: PESERTA + HRD + ADMIN
                // --------------------------------------------------------
                api.GET("/chat/tiket/:id/pesan", chatH.GetPesan)
                api.POST("/chat/tiket/:id/pesan", chatH.KirimPesan)
                api.GET("/chat/knowledge", chatH.GetKnowledge)

                // --------------------------------------------------------
                // HRD ROUTES
                // --------------------------------------------------------
                hrd := api.Group("/hrd")
                hrd.Use(middleware.RoleRequired(models.RoleHRD, models.RoleAdmin))
                {
                        // Pelaksanaan
                        hrd.POST("/pelaksanaan/pengajuan/:pengajuan_id", pelaksanaanH.SetJadwal)
                        hrd.GET("/pelaksanaan", pelaksanaanH.GetAll)
                        hrd.PATCH("/pelaksanaan/:id/nilai", pelaksanaanH.SetNilai)
                        hrd.POST("/pelaksanaan/:id/sertifikat", pelaksanaanH.GenerateSertifikat)

                        // Absensi
                        hrd.GET("/absensi/pelaksanaan/:id", absensiH.GetByPelaksanaan)
                        hrd.PATCH("/absensi/:id/approve", absensiH.Approve)

                        // Chat management
                        hrd.GET("/chat/tiket", chatH.GetAllTiket)
                        hrd.PATCH("/chat/tiket/:id/status", chatH.UpdateStatus)
                        hrd.POST("/chat/knowledge", chatH.UpsertKnowledge)

                        // Landing page management
                        hrd.PUT("/landing/content", landingH.UpdateContent)
                        hrd.GET("/landing/faq/all", landingH.GetFAQAll)
                        hrd.POST("/landing/faq", landingH.UpsertFAQ)
                        hrd.DELETE("/landing/faq/:id", landingH.DeleteFAQ)

                        // HRD list (untuk assign)
                        hrd.GET("/admin/hrd-list", adminH.GetHRDList)
                }

                // --------------------------------------------------------
                // ADMIN ONLY ROUTES
                // --------------------------------------------------------
                admin := api.Group("/admin")
                admin.Use(middleware.RoleRequired(models.RoleAdmin))
                {
                        admin.GET("/stats", adminH.GetStats)
                        admin.GET("/users", adminH.GetUsers)
                        admin.POST("/users", adminH.CreateUser)
                        admin.PATCH("/users/:id/toggle", adminH.ToggleUser)

                        // Periode magang (hanya admin)
                        admin.GET("/periode", landingH.GetAllPeriode)
                        admin.POST("/periode", landingH.CreatePeriode)
                        admin.PATCH("/periode/:id/aktif", landingH.SetPeriodeAktif)
                }
        }

        return r
}
