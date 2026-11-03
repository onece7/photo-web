package router

import (
	"log/slog"
	"photo-web/config"
	handlers2 "photo-web/handlers"
	middleware2 "photo-web/middleware"
	service2 "photo-web/service"

	"github.com/gin-gonic/gin"
)

// Setup 初始化路由
func Setup(
	cfg *config.Config,
	logger *slog.Logger,
	authService *service2.AuthService,
	generationService *service2.GenerationService,
	historyService *service2.HistoryService,
) *gin.Engine {
	// 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware2.Logger(logger))
	r.Use(middleware2.CORS())
	r.Use(gin.Recovery())

	// 限流器 (100 req/s, burst 200)
	rateLimiter := middleware2.NewRateLimiter(100, 200)
	r.Use(rateLimiter.Handler())

	// 创建 JWT 管理器
	jwtManager := middleware2.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	// 创建处理器
	authHandler := handlers2.NewAuthHandler(authService)
	genHandler := handlers2.NewGenerationHandler(generationService)
	historyHandler := handlers2.NewHistoryHandler(historyService, cfg.Upload.Dir)

	// 健康检查
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 静态文件
	r.GET("/api/v1/files/*filepath", handlers2.ServeGeneratedFile(cfg.Upload.Dir))

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开路由 — 认证
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", middleware2.AuthRequired(authService, jwtManager), authHandler.Logout)
		}

		// 公开路由 — 生成 (风格列表)
		v1.GET("/generate/styles", genHandler.GetStyles)

		// 需认证路由 — 生成
		gen := v1.Group("/generate")
		gen.Use(middleware2.AuthRequired(authService, jwtManager))
		{
			gen.POST("/simple", genHandler.SimpleGenerate)
			gen.POST("/pro", genHandler.ProGenerate)
			gen.GET("/status/:id", genHandler.GetGenerationStatus)
		}

		// 需认证路由 — 历史
		history := v1.Group("/history")
		history.Use(middleware2.AuthRequired(authService, jwtManager))
		{
			history.GET("", historyHandler.GetHistory)
			history.GET("/:id", historyHandler.GetHistoryDetail)
			history.DELETE("/:id", historyHandler.DeleteHistory)
		}
	}

	return r
}
