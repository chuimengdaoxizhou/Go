package main

import (
	"net/http"
	"rate_limiter/limiter"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 替换不同限流器即可
	rl := limiter.NewFixedWindow(5, time.Second)

	r.Use(func(c *gin.Context) {
		if !rl.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
