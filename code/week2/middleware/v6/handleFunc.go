package v6

// HandleFunc 定义业务逻辑函数类型
// Tips: 该类型应与http.HandlerFunc类型一致 此处只是暂时定义一下这个类型
type HandleFunc func(ctx *Context)
