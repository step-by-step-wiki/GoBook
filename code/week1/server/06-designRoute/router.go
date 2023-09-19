package designRoute

// router 路由森林 用于支持对路由树的操作
type router struct {
	// trees 路由森林 按HTTP动词组织路由树
	// 该map中 key为HTTP动词 value为路由树的根节点
	// 即: 每个HTTP动词对应一棵路由树 指向每棵路由树的根节点
	trees map[string]*node
}

// newRouter 创建路由森林
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// AddRoute 注册路由到路由森林中的路由树上
func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {
	// TODO: implement me
	panic("implement me")
}
