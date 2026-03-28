package main

import (
	"log"

	"github.com/ChenYujunjks/FlashCache/internal/cache"
	"github.com/ChenYujunjks/FlashCache/internal/handler"
	"github.com/ChenYujunjks/FlashCache/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	store := cache.NewInMemoryStore()
	cacheService := service.NewCacheService(store)
	cacheHandler := handler.NewCacheHandler(cacheService)

	r := gin.Default()

	handler.RegisterHealthRoutes(r)
	cacheHandler.RegisterRoutes(r)

	log.Println("flashcache server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
