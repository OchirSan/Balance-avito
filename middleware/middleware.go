package middleware

import (
	"avito/Balance-avito/prometheus"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// statusRecorder recodrs response status
type statusRecorder struct {
	http.ResponseWriter
	Status string
}

// WriteHeader reimplements WriteHeader method for ResponseWriter interface
func (rec *statusRecorder) WriteHeader(code int) {
	rec.Status = strconv.Itoa(code)
	rec.ResponseWriter.WriteHeader(code)
}

// Metrics middlware
func Metrics(p *prometheus.Prometheus) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := statusRecorder{ResponseWriter: w}
			now := time.Now()
			next.ServeHTTP(&recorder, r)
			path := mux.CurrentRoute(r).GetName()
			if path == "" {
				path = r.URL.Path
			}
			p.Reqs.WithLabelValues(r.Method, recorder.Status, path).Add(1)
			p.Latency.WithLabelValues(r.Method, recorder.Status, path).
				Observe(time.Since(now).Seconds())
		})
	}
}

