package middleware

import (
        "github.com/gin-contrib/cors"
        "github.com/gin-gonic/gin"
        "github.com/telpp/emagang/internal/config"
)

func CORS() gin.HandlerFunc {
        corsConfig := cors.Config{
                AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
                AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
                ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
                AllowCredentials: true,
        }
        if config.App.AppEnv == "development" {
                corsConfig.AllowAllOrigins = true
                corsConfig.AllowCredentials = false
        } else {
                corsConfig.AllowOrigins = []string{config.App.FrontendURL}
        }
        return cors.New(corsConfig)
}
