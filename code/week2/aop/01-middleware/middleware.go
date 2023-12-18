package middleware

// Middleware 中间件 用来包装 HandleFunc 返回一个新的 HandleFunc
// 这种入参和返回值均为一个函数的设计 是函数式编程 通过这种方式 可以将多个中间件串联起来
// 函数式的洋葱模式 或者叫 函数式的责任链模式
type Middleware func(next HandleFunc) HandleFunc

// Middleware 中间件接口
//type Middleware interface {
//	// Invoke 包装 HandleFunc 返回一个新的 HandleFunc
//	Invoke(next HandleFunc) HandleFunc
//}

// Interceptor 拦截器接口
//type Interceptor interface {
//	// Before 前置拦截器(在请求处理之前执行)
//	Before(ctx *Context)
//	// After 后置拦截器(在请求处理之后执行)
//	After(ctx *Context)
//	// Surround 环绕拦截器(在请求处理前后执行)
//	Surround(ctx *Context)
//}

// HandleChan 中间件链
//type HandleChan []HandleFunc

// HandlerChan 中间件链
//type HandlerChan struct {
//	// handlers 用于保存中间件链的切片
//	handlers []HandleFunc
//}

// Run 顺序执行中间件链上的每一个中间件
//func (h HandlerChan) Run(ctx *Context) {
//	for _, handler := range h.handlers {
//		handler(ctx)
//	}
//}

// HandleFuncNext 用于演示可控制是否执行下一个中间件的中间件函数
//type HandleFuncNext func(ctx *Context) (next bool)

// HandlerChan 中间件链
//type HandlerChan struct {
//	// handlers 用于保存中间件链的切片
//	handlers []HandleFuncNext
//}

// Run 顺序执行中间件链上的每一个中间件 直到某个中间件指定不再执行下一个中间件
//func (h HandlerChan) Run(ctx *Context) {
//	for _, handler := range h.handlers {
//		next := handler(ctx)
//		if !next {
//			return
//		}
//	}
//}

// Net 责任链的网状结构
//type Net struct {
//	// handlers 用于保存责任链的切片 注意切片中的每一个元素都是一条责任链
//	handlers []ConcurrentHandleFunc
//}

// Run 执行责任网上的每一个责任链
//func (n Net) Run(ctx *Context) {
//	wg := sync.WaitGroup{}
//	for _, handler := range n.handlers {
//		h := handler
//		if h.concurrent {
//			wg.Add(1)
//			go func() {
//				h.Run(ctx)
//				wg.Done()
//			}()
//		} else {
//			h.Run(ctx)
//		}
//	}
//	wg.Wait()
//}

// ConcurrentHandleFunc 允许并发执行责任链上的每一个中间件
//type ConcurrentHandleFunc struct {
//	// concurrent 标识是否允许并发执行的标量
//	concurrent bool
//	// handlers 用于保存中间件链的切片
//	handlers []*ConcurrentHandleFunc
//}

// Run 执行责任链上的每一个中间件 通过并发标识决定是否允许并发执行
//func (c ConcurrentHandleFunc) Run(ctx *Context) {
//	for _, handler := range c.handlers {
//		h := handler
//		if h.concurrent {
//			go h.Run(ctx)
//		} else {
//			h.Run(ctx)
//		}
//	}
//}
