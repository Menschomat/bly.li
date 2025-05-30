package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	instrumentHandlerOnce sync.Once
	httpRequestsTotal     = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "blyli_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "blyli_http_request_duration_seconds",
			Help:    "Request latency by handler",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method"},
	)
)

func SlogLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Debug("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration", time.Since(start),
				"remote", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		}
		return http.HandlerFunc(fn)
	}
}

func InstrumentHandler(next http.Handler) http.Handler {
	instrumentHandlerOnce.Do(func() {
		prometheus.MustRegister(httpRequestsTotal)
		prometheus.MustRegister(httpRequestDuration)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusResponseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()
		handler := r.URL.Path // or r.RoutePattern() with chi v5 and Go 1.22+
		httpRequestsTotal.WithLabelValues(handler, r.Method, http.StatusText(ww.status)).Inc()
		httpRequestDuration.WithLabelValues(handler, r.Method).Observe(duration)
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
