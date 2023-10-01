package internalhttp

import (
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func LoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Log(
						zap.ErrorLevel,
						"trace", zap.Any("Stack", debug.Stack()),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Log(
				zap.InfoLevel, "http logging",
				zap.Any("HOST", r.Host),
				zap.Any("UserAgent", r.UserAgent()),
				zap.Any("Status", wrapped.status),
				zap.Any("method", r.Method),
				zap.Any("path", r.URL.EscapedPath()),
				zap.Any("start", start),
				zap.Any("duration", time.Since(start)),
			)
		}

		return http.HandlerFunc(fn)
	}
}
