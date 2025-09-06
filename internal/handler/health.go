package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/insanjati/fitbyte/internal/cache"
	"github.com/insanjati/fitbyte/internal/database"
)

type HealthHandler struct {
	db    *database.DB
	cache *cache.Redis
}

func NewHealthHandler(db *database.DB, cache *cache.Redis) *HealthHandler {
	return &HealthHandler{
		db:    db,
		cache: cache,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	// Check database connection
	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "down",
			"database": "disconnected",
			"cache":    "unknown",
		})
		return
	}

	// Check cache connection
	ctx := c.Request.Context()
	if err := h.cache.SetExp(ctx, "health_check", "ok", 10*time.Second); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"database": "connected",
			"cache":    "disconnected",
		})
		return
	}

	// Clean up test key
	_ = h.cache.Delete(ctx, "health_check")

	c.JSON(http.StatusOK, gin.H{
		"status":   "up",
		"database": "connected",
		"cache":    "connected",
	})
}
