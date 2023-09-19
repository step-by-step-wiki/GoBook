package designRouteWithChildTree

// router 路由森林 用于支持对路由树的操作
type router struct {
	// trees 路由森林 按HTTP动词组织路由树
	// 该map中 key为HTTP动词 value为路由树
	// 即: 每个HTTP动词对应一棵路由树
	trees map[string]tree
}

// newRouter 创建路由森林
func newRouter() *router {
	return &router{
		// 此时还不知道树的结构 malloc一个空的map即可
		trees: map[string]tree{},
	}
}

// AddRoute 注册路由到路由森林中的路由树上
func (a *router) AddRoute(method string, path string, handleFunc HandleFunc) {
	// TODO: implement me
	panic("implement me")
}
