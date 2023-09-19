package designRouteWithChildTree

// node 路由树的节点
type node struct {
	// path 当前节点的路由路径
	path string

	// children 子路由路径到子节点的映射
	children map[string]*node

	// HandleFunc 路由对应的业务逻辑
	HandleFunc
}
