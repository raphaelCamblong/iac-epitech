package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimit returns a middleware that limits requests per IP.
// When limit exceeded, returns 429 Too Many Requests.
func RateLimit(requestsPerSecond float64, burst int) func(http.Handler) http.Handler {
	type entry struct {
		limiter *rate.Limiter
		last    time.Time
	}
	var (
		mu    sync.Mutex
		cache = make(map[string]*entry)
	)
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			mu.Lock()
			for k, v := range cache {
				if time.Since(v.last) > 2*time.Minute {
					delete(cache, k)
				}
			}
			mu.Unlock()
		}
	}()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				ip = fwd
			}
			mu.Lock()
			e, ok := cache[ip]
			if !ok {
				e = &entry{limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst), last: time.Now()}
				cache[ip] = e
			}
			e.last = time.Now()
			lim := e.limiter
			mu.Unlock()
			if !lim.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
