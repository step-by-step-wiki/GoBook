package access_log

import (
	"encoding/json"
	"web"
)

// MiddlewareBuilder accessLog中间件的构建器
type MiddlewareBuilder struct {
	logFunc func(logContent string) // logFunc 记录日志的函数,该函数用于提供给外部,让使用者自定义日志记录的方式
}

// SetLogFunc 设置记录日志的函数并返回构建器.返回构建器的目的在于:在使用时可以链式调用
// 即:使用时可以写出如下代码:
// m := &MiddlewareBuilder{}
//
//	m.SetLogFunc(func(logContent string) {
//		log.Print(logContent)
//	}).Build()
func (m *MiddlewareBuilder) SetLogFunc(logFunc func(logContent string)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

// Build 构建中间件
func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 调用业务处理函数完成后记录日志
			defer func() {
				accessLogObj := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}

				accessLogBytes, _ := json.Marshal(accessLogObj)
				m.logFunc(string(accessLogBytes))
			}()
			next(ctx)
		}
	}
}
