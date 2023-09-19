package tdd

import "strings"

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
	// step1. 找到路由树
	root, ok := r.trees[method]
	// 如果没有找到路由树,则创建一棵路由树
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// step2. 切割path
	// Tips: 去掉前导的"/" 否则直接切割出来的第一个元素为空字符串
	// Tips: 以下代码是老师写的去掉前导的"/"的方式 我认为表达力有点弱 但是性能应该会好于strings.TrimLeft
	// path = path[1:]
	path = strings.TrimLeft(path, "/")
	segments := strings.Split(path, "/")

	// step3. 为路由树添加路由
	// Tips: 此处我认为用target指代要添加路由的节点更好理解
	target := root
	for _, segment := range segments {
		// 如果路由树中途有节点没有创建,则创建该节点;
		// 如果路由树中途存在子节点,则找到该子节点
		child := target.childOrCreate(segment)
		// 继续为子节点创建子节点
		target = child
	}
	// 为目标节点设置HandleFunc
	target.HandleFunc = handleFunc
}
