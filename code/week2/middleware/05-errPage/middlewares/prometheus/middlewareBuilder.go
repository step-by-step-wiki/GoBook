package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
	"web"
)

// MiddlewareBuilder prometheus中间件构建器
type MiddlewareBuilder struct {
	Namespace string // Namespace APP名称
	Subsystem string // Subsystem 子系统/模块名称
	Name      string // Name 指标名称
	Help      string // Help 指标的描述信息
}

// Build 构建中间件
func (m *MiddlewareBuilder) Build() web.Middleware {
	labels := []string{
		"pattern", // pattern 命中的路由
		"method",  // method 请求方法
		"status",  // status 响应状态码
	}

	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name,
		Help:      m.Help,
		// 采样率 即百分比和误差
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.75: 0.01,
			0.90: 0.005,
			// 99线
			0.99: 0.001,
			// 999线
			0.999: 0.0001,
		},
	}, labels)

	// 注册指标
	prometheus.MustRegister(vector)

	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()

			// 为防止后续的中间件和HandleFunc中发生panic导致指标不记录
			// 故在此处将指标记录的代码放在defer中
			defer func() {
				// 响应时长
				duration := time.Since(startTime).Milliseconds()

				vector.WithLabelValues(
					// 命中的路由
					ctx.MatchRoute,
					// 请求方法
					ctx.Req.Method,
					// 响应状态码
					strconv.Itoa(ctx.RespStatusCode),
				).Observe(float64(duration))
			}()

			next(ctx)
		}
	}
}
