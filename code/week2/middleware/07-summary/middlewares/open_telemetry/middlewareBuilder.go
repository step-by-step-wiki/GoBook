package open_telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"web"
)

// instrumentationName 仪表盘名称 通常以包名作为仪表盘名称
// TODO: 如果真的把这个框架 作为一个独立的库发布 这里要改成github.com/xxx/xxx这样的形式
const instrumentationName = "web/middlewares/open_telemetry"

// MiddlewareBuilder openTelemetry中间件构建器
type MiddlewareBuilder struct {
	Tracer trace.Tracer // Tracer 追踪器
}

// Build 构建中间件
func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		if m.Tracer == nil {
			m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
		}

		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			carrier := propagation.HeaderCarrier(ctx.Req.Header)
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, carrier)

			reqCtx, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()

			attributes := []attribute.KeyValue{
				attribute.String("http.method", ctx.Req.Method),
				attribute.String("http.url", ctx.Req.URL.String()),
				attribute.String("http.scheme", ctx.Req.URL.Scheme),
				attribute.String("http.host", ctx.Req.Host),
			}
			span.SetAttributes(attributes...)

			ctx.Req = ctx.Req.WithContext(reqCtx)

			next(ctx)

			if ctx.MatchRoute != "" {
				span.SetName(ctx.MatchRoute)
			}

			// 强制类型转换(假定这里我们有一个自己实现的MyResponseWriter)
			// 问题在于:你根本不知道ctx.Resp是什么类型的 可能是http包内的私有类型
			// 还有可能是使用者通过中间件提供的自定义类型
			// myResp, ok := ctx.Resp.(MyResponseWriter)

			// 4. 请求完成后记录响应码
			span.SetAttributes(attribute.Int("http.status_code", ctx.RespStatusCode))
		}
	}
}
