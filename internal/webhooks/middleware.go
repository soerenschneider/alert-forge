package webhooks

import (
	"net/http"
	"strconv"

	"github.com/soerenschneider/alert-forge/internal/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wr, r)
		metrics.Requests.WithLabelValues(r.URL.Path, strconv.Itoa(wr.statusCode)).Inc()
	})
}
