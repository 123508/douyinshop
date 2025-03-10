package util

import (
	"context"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"time"
)

// 定义请求计数器
var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kitex_requests_total",
			Help: "Total number of Kitex requests received",
		},
		[]string{"method", "status"},
	)

	// 定义响应时间直方图
	responseDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "kitex_response_duration_seconds",
			Help: "Kitex response duration in seconds",
		},
		[]string{"method"},
	)
	// 新增：当前并发量
	concurrentRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kitex_concurrent_requests",
			Help: "Current number of concurrent Kitex requests",
		},
		[]string{"method"},
	)

	// 新增：RPS 计算依赖的计数器
	requestsPerSecond = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kitex_requests_per_second",
			Help: "Total requests for RPS calculation",
		},
		[]string{"method"},
	)
)

// 初始化 logrus 配置（可选，例如 JSON 格式）
func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// logrus.SetLevel(logrus.InfoLevel) // 设置日志级别
}

func MonitorMiddleware() endpoint.Middleware {

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			startTime := time.Now()
			method := rpcinfo.GetRPCInfo(ctx).To().Method()

			// 记录并发量 +1（请求开始时）
			concurrentRequests.WithLabelValues(method).Inc()
			// 确保请求结束时减少并发量
			defer concurrentRequests.WithLabelValues(method).Dec()

			// 记录 RPS 计数器 +1（每次请求计数）
			requestsPerSecond.WithLabelValues(method).Inc()

			// 创建带有固定字段的 logrus.Entry
			logEntry := logrus.WithFields(logrus.Fields{
				"component": "Prometheus", // 标识日志来源
				"method":    method,
			})

			// 记录请求开始和参数
			logEntry.WithField("request", req).Info("开始处理请求")

			// 调用实际处理函数
			err = next(ctx, req, resp)

			duration := time.Since(startTime)
			statusLabel := "success"
			if err != nil {
				statusLabel = "failure"
				// 记录错误（附带错误信息和耗时）
				logEntry.WithError(err).
					WithField("duration", duration.String()).
					Error("请求失败")
			} else {
				// 记录成功（附带耗时）
				logEntry.WithField("duration", duration.String()).
					Info("请求成功")
			}

			// 记录响应结果
			logEntry.WithField("response", resp).Debug("响应结果") // 可根据级别调整

			// 记录 Prometheus 指标
			requestCounter.WithLabelValues(method, statusLabel).Inc()
			responseDuration.WithLabelValues(method).Observe(duration.Seconds())

			// 记录指标提交状态
			logEntry.WithFields(logrus.Fields{
				"status":           statusLabel,
				"duration_seconds": duration.Seconds(),
			}).Info("指标已记录")

			return err
		}
	}
}
