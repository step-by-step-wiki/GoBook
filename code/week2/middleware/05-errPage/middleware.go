package web

// Middleware 中间件 用来包装 HandleFunc 返回一个新的 HandleFunc
// 这种入参和返回值均为一个函数的设计 是函数式编程 通过这种方式 可以将多个中间件串联起来
// 函数式的洋葱模式 或者叫 函数式的责任链模式
type Middleware func(next HandleFunc) HandleFunc
