package handler

import (
	"errors"
	"net/http"

	"github.com/ChenYujunjks/FlashCache/internal/model"
	"github.com/ChenYujunjks/FlashCache/internal/service"

	"github.com/gin-gonic/gin"
)

type CacheHandler struct {
	cacheService *service.CacheService
}

func NewCacheHandler(cacheService *service.CacheService) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
	}
}

func (h *CacheHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.PUT("/cache/:key", h.Set)
		v1.GET("/cache/:key", h.Get)
		v1.DELETE("/cache/:key", h.Delete)
	}
}

func (h *CacheHandler) Set(c *gin.Context) {
	key := c.Param("key")

	var req model.SetCacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	if err := h.cacheService.Set(key, req.Value, req.TTLSeconds); err != nil {
		status := http.StatusBadRequest
		c.JSON(status, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: gin.H{
			"key":         key,
			"value":       req.Value,
			"ttl_seconds": req.TTLSeconds,
		},
	})
}

func (h *CacheHandler) Get(c *gin.Context) {
	key := c.Param("key")

	value, err := h.cacheService.Get(key)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: gin.H{
			"key":   key,
			"value": value,
		},
	})
}

func (h *CacheHandler) Delete(c *gin.Context) {
	key := c.Param("key")

	err := h.cacheService.Delete(key)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: gin.H{
			"key":     key,
			"deleted": true,
		},
	})
}
