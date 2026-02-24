package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

const requestIDHeader = "X-Request-ID"

type APIMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewAPIMetrics(registerer prometheus.Registerer) *APIMetrics {
	m := &APIMetrics{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests.",
			},
			[]string{"method", "route", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "api",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "Request duration in seconds.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "route", "status"},
		),
	}

	registerer.MustRegister(m.requestsTotal, m.requestDuration)
	return m
}

func (m *APIMetrics) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		statusCode := strconv.Itoa(c.Writer.Status())
		m.requestsTotal.WithLabelValues(c.Request.Method, route, statusCode).Inc()
		m.requestDuration.WithLabelValues(c.Request.Method, route, statusCode).Observe(time.Since(start).Seconds())
	}
}

func RequestLogger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Writer.Header().Set(requestIDHeader, requestID)

		start := time.Now()
		c.Next()

		logger.Info().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Msg("http request")
	}
}

func Recoverer(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error().Any("panic", rec).Bytes("stack", debug.Stack()).Msg("panic recovered")
				c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			}
		}()

		c.Next()
	}
}

func newRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(bytes)
}
