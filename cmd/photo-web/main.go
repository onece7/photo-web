package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"photo-web/internal/bootstrap"
	"photo-web/internal/migration"
	"syscall"
	"time"

	"photo-web/internal/config"
	"photo-web/internal/repository"
	"photo-web/internal/router"
	"photo-web/internal/service"
	"photo-web/pkg/redis"
	"photo-web/pkg/response"
)

func main() {
	configPath := flag.String("config", "config.yaml", "")
	flag.Parse()
	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	db, err := bootstrap.InitDB(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库
	migration.Migrate(db)

	// 生成种子风格
	migration.SeedStyles(db)

	// 初始化 Redis
	rdb, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}

	// 初始化统一响应
	response.Init()

	// 依赖注入 - Repository 层
	userRepo := repository.NewUserRepository(db)
	styleRepo := repository.NewStyleRepository(db)
	generationRepo := repository.NewGenerationRepository(db)

	// 依赖注入 - Service 层
	authService := service.NewAuthService(userRepo, rdb, cfg)
	sdClient := service.NewSDClient(cfg)
	generationService := service.NewGenerationService(generationRepo, styleRepo, sdClient, cfg)
	historyService := service.NewHistoryService(generationRepo)

	// 初始化路由
	r := router.Setup(cfg, authService, generationService, historyService)

	// 启动服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 优雅关闭
	go func() {
		fmt.Printf("服务启动: http://localhost:%d\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务关闭: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}
	fmt.Println("服务已安全关闭")
}

