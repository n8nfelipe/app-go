package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	rateLimitPerSecond = rate.Limit(0.5) // 1 requisição a cada 2s
	rateLimitBurst     = 5
	visitorTTL         = 10 * time.Minute
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitorsMu sync.Mutex
	visitors   = make(map[string]*visitor)
)

func init() {
	go cleanupVisitors()
}

func visitorFor(key string) *rate.Limiter {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	v, ok := visitors[key]
	if !ok {
		v = &visitor{limiter: rate.NewLimiter(rateLimitPerSecond, rateLimitBurst)}
		visitors[key] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for range time.Tick(time.Minute) {
		visitorsMu.Lock()
		for key, v := range visitors {
			if time.Since(v.lastSeen) > visitorTTL {
				delete(visitors, key)
			}
		}
		visitorsMu.Unlock()
	}
}

// RateLimit limita requisições por IP em endpoints sensíveis (login, cadastro).
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !visitorFor(c.ClientIP()).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Muitas requisições, tente novamente mais tarde"})
			c.Abort()
			return
		}
		c.Next()
	}
}
