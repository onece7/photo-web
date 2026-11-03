package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	bootstrap2 "photo-web/bootstrap"
	"photo-web/config"
	middleware2 "photo-web/middleware"
	repository2 "photo-web/repository"
	"photo-web/router"
	service2 "photo-web/service"
	"syscall"
	"time"
)

func main() {
	configPath := flag.String("config", "config.yaml", "")
	flag.Parse()
	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	logger, err := bootstrap2.NewLogger(cfg, os.Stdout)
	if err != nil {
		log.Fatalf("日志初始化失败：%v", err)
	}
	slog.SetDefault(logger)

	// 初始化数据库
	db, err := bootstrap2.InitDB(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库
	err = repository2.Migrate(db)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 生成种子风格
	repository2.SeedStyles(db, logger)

	// 初始化 Redis
	rdb, err := middleware2.NewClient(cfg)
	if err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}

	// 依赖注入 - Repository 层
	userRepo := repository2.NewUserRepository(db)
	styleRepo := repository2.NewStyleRepository(db)
	generationRepo := repository2.NewGenerationRepository(db)

	// 依赖注入 - Service 层
	authService := service2.NewAuthService(userRepo, rdb, cfg)
	sdClient := service2.NewSDClient(cfg)
	generationService := service2.NewGenerationService(generationRepo, styleRepo, sdClient, cfg, logger)
	historyService := service2.NewHistoryService(generationRepo)

	// 初始化路由
	r := router.Setup(cfg, logger, authService, generationService, historyService)

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
