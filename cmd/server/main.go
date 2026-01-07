package main

import (
	"context"
	"log"
	"notification-system/internal/api"
	"notification-system/internal/config"
	"notification-system/internal/deliver"
	"notification-system/internal/model"
	"notification-system/internal/platform"
	"notification-system/internal/platform/impl"
	"notification-system/internal/repository"
	"notification-system/internal/semantic"
	"notification-system/internal/service"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()

	var repo repository.Repository
	if cfg.UseMySQL {
		log.Println("Using MySQL repository...")
		mysqlRepo, err := repository.NewMySQLRepository(cfg.DB.DSN)
		if err != nil {
			log.Fatalf("Failed to connect to MySQL: %v", err)
		}
		repo = mysqlRepo
		log.Println("MySQL connected successfully")
	} else {
		log.Println("Using in-memory repository...")
		repo = repository.NewInMemoryRepository()
	}

	platformMgr := platform.NewPlatformManager()

	platformMgr.RegisterPlatform(impl.NewHTTP2_Platform())
	platformMgr.RegisterPlatform(impl.NewHTTPPlatform())
	platformMgr.RegisterPlatform(impl.NewEmailPlatform())
	platformMgr.RegisterPlatform(impl.NewSMSPlatform())

	semanticMgr := semantic.NewSemanticManager()
	semanticMgr.RegisterHandler(string(model.DeliverySemanticAtLeastOnce), semantic.NewAtLeastOnceHandler())
	semanticMgr.RegisterHandler(string(model.DeliverySemanticAtMostOnce), semantic.NewAtMostOnceHandler())
	semanticMgr.RegisterHandler(string(model.DeliverySemanticExactlyOnce), semantic.NewExactlyOnceHandler(repo))

	deliverer := deliver.NewDeliverer(repo, platformMgr, semanticMgr)

	svc := service.NewService(repo, deliverer, semanticMgr)

	handler := api.NewHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc.Start(ctx)

	r := gin.Default()
	handler.RegisterRoutes(r)

	go func() {
		if err := r.Run(":" + strconv.Itoa(cfg.Server.Port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("Server stopped")
}
