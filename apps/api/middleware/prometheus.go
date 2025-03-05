package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"time"
)

// 定义 Prometheus 指标
var (
	// HTTP 请求总数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP 请求延迟
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	// HTTP 请求大小
	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	// HTTP 响应大小
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	// 当前活跃请求数
	httpActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Current number of active HTTP requests",
		},
		[]string{"method", "path"},
	)

	// 认证失败次数
	authFailuresTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_failures_total",
			Help: "Total number of authentication failures",
		},
		[]string{"method", "path"},
	)

	// 业务错误次数
	businessErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_errors_total",
			Help: "Total number of business logic errors",
		},
		[]string{"method", "path", "error_type"},
	)
)

// 初始化日志配置
func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// PrometheusMiddleware 创建 Prometheus 监控中间件
func PrometheusMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()
		method := string(c.Method())
		path := c.FullPath()
		if path == "" {
			path = string(c.Request.URI().Path())
		}

		// 记录活跃请求数
		httpActiveRequests.WithLabelValues(method, path).Inc()
		defer httpActiveRequests.WithLabelValues(method, path).Dec()

		// 创建日志条目
		logEntry := logrus.WithFields(logrus.Fields{
			"method": method,
			"path":   path,
			"ip":     c.ClientIP(),
		})

		// 记录请求开始
		logEntry.Info("开始处理请求")

		// 记录请求大小
		requestSize := len(c.Request.Body())
		httpRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))

		// 处理请求
		c.Next(ctx)

		// 计算处理时间
		duration := time.Since(startTime)
		status := c.Response.StatusCode()

		// 记录响应大小
		responseSize := len(c.Response.Body())
		httpResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))

		// 记录请求总数和延迟
		httpRequestsTotal.WithLabelValues(method, path, getStatusLabel(status)).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

		// 记录错误
		if status >= 400 {
			if status == 401 || status == 403 {
				authFailuresTotal.WithLabelValues(method, path).Inc()
			} else if status >= 500 {
				businessErrorsTotal.WithLabelValues(method, path, "server_error").Inc()
			} else {
				businessErrorsTotal.WithLabelValues(method, path, "client_error").Inc()
			}
		}

		// 记录请求完成
		logEntry.WithFields(logrus.Fields{
			"status":        status,
			"duration":      duration.String(),
			"request_size":  requestSize,
			"response_size": responseSize,
		}).Info("请求处理完成")
	}
}

// getStatusLabel 将 HTTP 状态码转换为 Prometheus 标签
func getStatusLabel(status int) string {
	switch {
	case status < 200:
		return "invalid"
	case status < 300:
		return "success"
	case status < 400:
		return "redirect"
	case status < 500:
		return "client_error"
	default:
		return "server_error"
	}
}

// RegisterPrometheus 注册 Prometheus 中间件
func RegisterPrometheus(h *server.Hertz) {
	// 添加 Prometheus 中间件
	h.Use(PrometheusMiddleware())
}
