package router

import (
        "context"
        "net/http"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/config"
        "github.com/telpp/emagang/internal/database"
        "github.com/telpp/emagang/internal/handler"
        "github.com/telpp/emagang/internal/middleware"
        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
        "github.com/telpp/emagang/internal/routes"
        "github.com/telpp/emagang/internal/service"
        ws "github.com/telpp/emagang/internal/websocket"
)

func Setup() *gin.Engine {
        r := gin.Default()

        r.Use(middleware.CORS())
        r.MaxMultipartMemory = 110 << 20

        // Repositories
        userRepo         := repository.NewUserRepository(database.DB)
        pengajuanRepo    := repository.NewPengajuanRepository(database.DB)
        dokumenRepo      := repository.NewDokumenRepository(database.DB)
        pelaksanaanRepo  := repository.NewPelaksanaanRepository(database.DB)
        absensiRepo      := repository.NewAbsensiRepository(database.DB)
        absensiConfigRepo := repository.NewAbsensiConfigRepository(database.DB)
        izinSakitRepo    := repository.NewIzinSakitRepository(database.DB)
        chatRepo         := repository.NewChatRepository(database.DB)
        notifRepo        := repository.NewNotifikasiRepository(database.DB)
        landingRepo      := repository.NewLandingRepository(database.DB)
        divisiRepo       := repository.NewDivisiRepository(database.DB)
        laporanRepo      := repository.NewLaporanRepository(database.DB)
        penilaianRepo    := repository.NewPenilaianRepository(database.DB)
        perpanjanganRepo := repository.NewPerpanjanganRepository(database.DB)
        knowledgeRepo    := repository.NewKnowledgeRepository(database.DB)

        repository.SetDB(database.DB)

        // Services
        hub           := ws.GlobalHub
        fcmSvc        := service.NewFCMService(config.App.FirebaseProjectID, config.App.FirebaseServiceAccountJSON)
        notifSvc      := service.NewNotifikasiService(notifRepo, hub, fcmSvc)
        authSvc       := service.NewAuthService(userRepo)
        emailSvc      := service.NewEmailService()
        pengajuanSvc  := service.NewPengajuanService(pengajuanRepo, notifSvc, userRepo, emailSvc, dokumenRepo)
        chatSvc       := service.NewChatService(chatRepo, notifSvc, userRepo, hub)
        nlpSvc        := service.NewNLPService(knowledgeRepo)
        sertifikatSvc := service.NewSertifikatService(pelaksanaanRepo, pengajuanRepo, notifSvc, emailSvc)
        cleanupSvc       := service.NewCleanupService(userRepo, emailSvc)
        absensiReminderSvc := service.NewAbsensiReminderService(absensiRepo, absensiConfigRepo, notifSvc)

        // Handlers
        authH         := handler.NewAuthHandler(authSvc)
        pengajuanH    := handler.NewPengajuanHandler(pengajuanSvc)
        dokumenH      := handler.NewDokumenHandler(dokumenRepo, pengajuanRepo, notifSvc)
        pelaksanaanH  := handler.NewPelaksanaanHandler(pelaksanaanRepo, pengajuanRepo, sertifikatSvc)
        absensiH      := handler.NewAbsensiHandler(absensiRepo, pelaksanaanRepo, pengajuanRepo, absensiConfigRepo, divisiRepo)
        absensiConfigH := handler.NewAbsensiConfigHandler(absensiConfigRepo)
        izinSakitH    := handler.NewIzinSakitHandler(izinSakitRepo, pelaksanaanRepo)
        chatH         := handler.NewChatHandler(chatSvc)
        notifH        := handler.NewNotifikasiHandler(notifSvc, chatRepo)
        landingH      := handler.NewLandingHandler(landingRepo)
        adminH        := handler.NewAdminHandler(userRepo, dokumenRepo, emailSvc)
        divisiH       := handler.NewDivisiHandler(divisiRepo)
        laporanH      := handler.NewLaporanHandler(laporanRepo, pelaksanaanRepo)
        penilaianH    := handler.NewPenilaianHandler(penilaianRepo)
        perpanjanganH := handler.NewPerpanjanganHandler(perpanjanganRepo, pelaksanaanRepo, notifSvc, emailSvc)
        badgeH        := handler.NewBadgeHandler()
        knowledgeH    := handler.NewKnowledgeHandler(knowledgeRepo, nlpSvc)

        r.Static("/uploads", config.App.UploadDir)

        r.GET("/api/health", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"status": "ok", "app": "e-Magang TELPP"})
        })

        // Route groups
        public  := r.Group("/api")
        api     := r.Group("/api")
        api.Use(middleware.AuthRequired())
        peserta := api.Group("")
        peserta.Use(middleware.RoleRequired(models.RolePeserta))
        hrd := api.Group("")
        hrd.Use(middleware.RoleRequired(models.RoleHRD, models.RoleAdmin))
        admin := api.Group("/admin")
        admin.Use(middleware.AuthRequired(), middleware.RoleRequired(models.RoleAdmin))

        // WebSocket
        public.GET("/ws", ws.ServeWS(hub))

        // Register routes
        routes.RegisterAuthRoutes(public, api, authH)
        routes.RegisterPengajuanRoutes(public, peserta, api, hrd, pengajuanH)
        routes.RegisterDokumenRoutes(public, peserta, api, hrd, dokumenH)
        routes.RegisterPelaksanaanRoutes(peserta, hrd, api,pelaksanaanH)
        routes.RegisterAbsensiRoutes(peserta, hrd, absensiH)
        routes.RegisterAbsensiConfigRoutes(api, admin, absensiConfigH)
        routes.RegisterIzinSakitRoutes(peserta, hrd, izinSakitH)
        routes.RegisterChatRoutes(peserta, api, hrd, chatH)
        routes.RegisterNotifikasiRoutes(api, notifH)
        routes.RegisterLandingRoutes(public, hrd, admin, landingH)
        routes.RegisterAdminRoutes(hrd, admin, adminH, dokumenH)
        routes.RegisterDivisiRoutes(api, admin, divisiH)
        routes.RegisterLaporanRoutes(peserta, hrd, api, laporanH)
        routes.RegisterPenilaianRoutes(api, hrd, penilaianH)
        routes.RegisterPerpanjanganRoutes(peserta, hrd, perpanjanganH)
        routes.RegisterBadgeRoutes(peserta, hrd, badgeH)
        routes.RegisterKnowledgeRoutes(public, admin, knowledgeH)

        // Auto-update status pelaksanaan setiap jam
        go func() {
                pelaksanaanRepo.AutoUpdateStatuses(context.Background())
                ticker := time.NewTicker(1 * time.Hour)
                defer ticker.Stop()
                for range ticker.C {
                        pelaksanaanRepo.AutoUpdateStatuses(context.Background())
                }
        }()

        // Timer hangus tiket — cek tiap 15 menit
        chatSvc.StartTimerGoroutine(context.Background())

        // Hapus otomatis akun peserta yang sudah selesai magang > 30 hari
        cleanupSvc.Start(context.Background())

        // Reminder absen masuk & pulang via push notif (cek tiap menit, WIB)
        absensiReminderSvc.Start(context.Background())

        return r
}
