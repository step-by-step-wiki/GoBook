package err_page

import "web"

// MiddlewareBuilder 错误页面中间件构造器
type MiddlewareBuilder struct {
	respPages map[int][]byte // respPages 用于存储响应码与其对应的错误页面 其中key为响应码 value为错误页面的内容
}

// Build 构造错误页面中间件
func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)

			// 判断响应码是否为需要篡改响应的响应码 如果是则篡改响应
			respPage, ok := m.respPages[ctx.RespStatusCode]
			if ok {
				ctx.RespData = respPage
			}
		}
	}
}

// AddCode 添加响应码与其对应的错误页面
func (m *MiddlewareBuilder) AddCode(status int, page []byte) *MiddlewareBuilder {
	if m.respPages == nil {
		m.respPages = make(map[int][]byte)
	}
	m.respPages[status] = page

	// Tips: 此处返回 *MiddlewareBuilder 是为了支持链式调用
	return m
}

// NewMiddlewareBuilder 初始化错误页面中间件构造器
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		respPages: make(map[int][]byte),
	}
}
