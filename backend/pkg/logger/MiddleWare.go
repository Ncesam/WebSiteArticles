package logger

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RequestCounter struct {
	LastTimeRequest time.Time
	RequestCount    int
}

var (
	clientResponse = make(map[string]*RequestCounter)
	mu             sync.RWMutex
)

func MiddleWare(logger *zap.Logger, maxRequests int, timeWindow time.Duration) gin.HandlerFunc {
	go startCleanupRoutine(5 * time.Minute)
	return func(c *gin.Context) {

		clientIp := c.ClientIP()

		clientRequestsInfo, exists := clientResponse[clientIp]
		if !exists {
			clientRequestsInfo = &RequestCounter{
				LastTimeRequest: time.Now(),
				RequestCount:    0,
			}
			clientResponse[clientIp] = clientRequestsInfo
		}

		if time.Since(clientRequestsInfo.LastTimeRequest) > timeWindow {
			clientRequestsInfo.RequestCount = 0
		}

		clientRequestsInfo.RequestCount++
		clientRequestsInfo.LastTimeRequest = time.Now()

		logger.Info("Request received", zap.String("ip", clientIp), zap.String("method", c.Request.Method), zap.String("url", c.Request.URL.String()))
		if clientRequestsInfo.RequestCount > maxRequests {
			logger.Warn("Too many requests", zap.String("ip", clientIp), zap.Int("count", clientRequestsInfo.RequestCount))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			return
		}
		c.Next()
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("GIN error", zap.Error(err))
			}
		}
	}

}

func startCleanupRoutine(ttl time.Duration) {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		now := time.Now()
		mu.Lock()
		for ip, info := range clientResponse {
			if now.Sub(info.LastTimeRequest) > ttl {
				delete(clientResponse, ip)
			}
		}
		mu.Unlock()
	}
}
