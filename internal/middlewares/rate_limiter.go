package middlewares

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/username/project-name/internal/config"
	"github.com/username/project-name/internal/errors"
	"github.com/username/project-name/internal/response"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*visitor
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
	expiry   time.Duration
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*visitor),
		rate:     rate.Every(cfg.Window / time.Duration(cfg.Requests)),
		burst:    cfg.Requests,
		expiry:   cfg.Window,
	}

	go rl.cleanup()
	return rl
}

func (r *RateLimiter) getLimiter(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	v, exists := r.limiters[ip]
	if !exists {
		limiter := rate.NewLimiter(r.rate, r.burst)
		r.limiters[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (r *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		for ip, v := range r.limiters {
			if time.Since(v.lastSeen) > r.expiry {
				delete(r.limiters, ip)
			}
		}
		r.mu.Unlock()
	}
}

func (r *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		limiter := r.getLimiter(getIP(req))
		if !limiter.Allow() {
			response.Error(w, req, errors.New("rate_limited", "rate limit exceeded", http.StatusTooManyRequests))
			return
		}

		next.ServeHTTP(w, req)
	})
}

func getIP(r *http.Request) string {
	forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if forwarded != "" {
		return forwarded
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return r.RemoteAddr
}
