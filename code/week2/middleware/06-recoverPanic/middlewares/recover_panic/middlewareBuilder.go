package recover_panic

import "web"

// MiddlewareBuilder 捕获panic中间件的构建器
type MiddlewareBuilder struct {
	StatusCode int                    // StatusCode 捕获panic时的响应状态码
	Data       []byte                 // Data 捕获panic时的响应数据
	LogFunc    func(ctx *web.Context) // LogFunc 捕获panic时的日志记录函数 (记录整个ctx)
	// LogFunc    func(err any)          // LogFunc 捕获panic时的日志记录函数 (记录panic的内容)
	// LogFunc    func(stack string)     // LogFunc 捕获panic时的日志记录函数 (记录调用栈)
}

// Build 构建捕获panic中间件
func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 捕获panic 篡改响应 并记录日志
			defer func() {
				// Tips: 这里的err的类型是ary，而不是error
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = m.Data

					// 记录日志
					if m.LogFunc != nil {
						m.LogFunc(ctx)
					}
				}
			}()

			next(ctx)
		}
	}
}
