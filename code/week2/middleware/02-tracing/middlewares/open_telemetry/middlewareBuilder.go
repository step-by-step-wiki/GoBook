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
			// step3. 尝试与客户端的span建立父子关系
			// 1. 从请求头中获取traceId和spanId (客户端的traceId和spanId是放在HTTP请求头中的)
			reqCtx := ctx.Req.Context()
			// Tips: 其实这个HeaderCarrier就是一个map[string][]string 跟request.Header是一样的
			carrier := propagation.HeaderCarrier(ctx.Req.Header)
			// 2. 基于客户端的traceId和spanId 重新封装成context.Context
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, carrier)

			// step1. 创建span
			// 1. Tracer.Start()方法的第1个参数是一个 context.Context 接口的实现
			// 而我们的 web.Context 并没有实现 context.Context 接口 (其实在这里拿到的是 context.Background() )
			// 2. Tracer.Start()方法的第2个参数是一个 span 的名称 通常以请求命中的路由作为 span 的名称
			// 但此时还不知道这个请求是否能命中路由 所以使用 "unknown" 作为 span 的名称
			// Tips: 此时拿到的就是基于客户端的traceId和spanId封装过的context.Context
			_, span := m.Tracer.Start(reqCtx, "unknown")
			// 4. 调用完成后要关闭span
			defer span.End()

			// step2. 设置span的标签(就是记录数据)
			attributes := []attribute.KeyValue{
				// 请求的HTTP动词
				attribute.String("http.method", ctx.Req.Method),
				// 请求的url
				attribute.String("http.url", ctx.Req.URL.String()),
				// 请求的scheme (http/https)
				attribute.String("http.scheme", ctx.Req.URL.Scheme),
				// 请求的host
				attribute.String("http.host", ctx.Req.URL.Host),
			}
			span.SetAttributes(attributes...)

			next(ctx)

			// 3. 只有当请求命中路由(实际上在这里已经执行完了对应的 web.HandleFunc 了)后
			// 才能确定请求命中的路由 因此在执行完 HandleFunc 后 再把span的名称改成请求命中的路由
			if ctx.MatchRoute != "" {
				span.SetName(ctx.MatchRoute)
			}
		}
	}
}

//// MiddlewareBuilder openTelemetry中间件构建器
//type MiddlewareBuilder struct {
//	tracer trace.Tracer // tracer 追踪器
//}
//
//// NewMiddlewareBuilder 创建中间件构建器
//func NewMiddlewareBuilder(tracer trace.Tracer) *MiddlewareBuilder {
//	return &MiddlewareBuilder{
//		tracer: tracer,
//	}
//}
