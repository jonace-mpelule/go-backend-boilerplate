package middlewares

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*visitor
	mu       sync.Mutex

	rate   rate.Limit
	burst  int
	expiry time.Duration
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(
	requests int,
	duration time.Duration,
) *RateLimiter {

	if requests <= 0 {
		requests = 10
	}

	if duration <= 0 {
		duration = time.Minute
	}

	rl := &RateLimiter{
		limiters: make(map[string]*visitor),
		rate:     rate.Every(duration / time.Duration(requests)),
		burst:    requests,
		expiry:   duration,
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
	for {
		time.Sleep(time.Minute)

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
		ip := getIP(req)

		limiter := r.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(
				w,
				"rate limit exceeded",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")

	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Real-IP")

	if ip != "" {
		return ip
	}

	return r.RemoteAddr
}
