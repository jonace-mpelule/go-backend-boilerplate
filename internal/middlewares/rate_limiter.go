package middlewares

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)

var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {

	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]

	if !exists {
		limiter = rate.NewLimiter(
			5,
			10,
		)
		visitors[ip] = limiter
	}

	return limiter
}

func RateLimit(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		ip := r.RemoteAddr

		limiter := getVisitor(ip)

		if !limiter.Allow() {
			http.Error(
				w,
				"rate limit exceeded",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
