package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vcsfrl/xm/internal/config"
	"golang.org/x/time/rate"
	"net/http"
)

// RateLimiter Middleware to check the rate limit.
func RateLimiter(config *config.Config) func(c *gin.Context) {
	// Define a rate limiter
	var limiter = rate.NewLimiter(rate.Limit(config.RateLimit), config.RateBurst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
