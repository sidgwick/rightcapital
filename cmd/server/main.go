package main

import (
	"context"
	"log"
	"notification-system/internal/api"
	"notification-system/internal/config"
	"notification-system/internal/deliver"
	"notification-system/internal/model"
	"notification-system/internal/platform"
	"notification-system/internal/repository"
	"notification-system/internal/semantic"
	"notification-system/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type DefaultPlatform struct{}

func (p *DefaultPlatform) Name() string {
	return "default"
}

func (p *DefaultPlatform) Deliver(ctx context.Context, target string, content interface{}) error {
	log.Printf("Delivering to %s: %v", target, content)
	return nil
}

func main() {
	cfg := config.NewConfig()

	repo := repository.NewInMemoryRepository()

	platformMgr := platform.NewPlatformManager()
	platformMgr.RegisterPlatform(&DefaultPlatform{})

	semanticMgr := semantic.NewSemanticManager()
	semanticMgr.RegisterHandler(string(model.DeliverySemanticAtLeastOnce), semantic.NewAtLeastOnceHandler())
	semanticMgr.RegisterHandler(string(model.DeliverySemanticAtMostOnce), semantic.NewAtMostOnceHandler())
	semanticMgr.RegisterHandler(string(model.DeliverySemanticExactlyOnce), semantic.NewExactlyOnceHandler())

	deliverer := deliver.NewDeliverer(repo, platformMgr, semanticMgr)

	svc := service.NewService(repo, deliverer, semanticMgr)

	handler := api.NewHandler(svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc.Start(ctx)

	r := gin.Default()
	handler.RegisterRoutes(r)

	go func() {
		if err := r.Run(":" + string(rune(cfg.Server.Port))); err != nil {
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
