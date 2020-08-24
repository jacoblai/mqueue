package limit

import (
	"golang.org/x/time/rate"
	"net/http"
)

var limiter = rate.NewLimiter(100000, 100000)

func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
