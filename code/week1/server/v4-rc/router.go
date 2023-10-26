package v4_rc

import (
	"fmt"
	"strings"
)

// router 路由森林
type router struct {
	// trees 路由森林 key为HTTP动词 value为HTTP对应路由树的根节点
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// addRoute 添加路由
func (r *router) addRoute(method string, path string, handle HandleFunc) {
	msg, ok := r.checkPath(path)
	if !ok {
		panic(msg)
	}

	// step1. 查找路由树,不存在则创建
	root, exist := r.trees[method]
	if !exist {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// step2. 在根节点上查找子节点 不存在则创建
	// step2.1 由于按/切割后 第一个元素为"" 也就是说如果传入的path为"/" 需要特殊处理
	if path == "/" {
		if root.HandleFunc != nil {
			msg = fmt.Sprintf("web: 路由冲突,重复注册路由 [%s]", path)
			panic(msg)
		}
		root.HandleFunc = handle
		return
	}

	// step2.2 从根节点开始 逐层查找
	target := root
	path = strings.TrimLeft(path, "/")
	pathSegments := strings.Split(path, "/")
	for _, pathSegment := range pathSegments {
		// 在当前节点上查找子节点
		child := target.findOrCreate(pathSegment)
		target = child
	}

	// 为目标节点创建HandleFunc
	if target.HandleFunc != nil {
		msg = fmt.Sprintf("web: 路由冲突,重复注册路由 [%s]", path)
		panic(msg)
	}
	target.HandleFunc = handle
}

// checkPath 检测路由是否合法
// 此处没有返回error 是因为设计上如果路由不合法 直接panic而非报错
// 所以此方法只返回 表示是否合法的标量以及表示不合法原因的字符串即可
func (r *router) checkPath(path string) (msg string, ok bool) {
	if path == "" {
		return "web: 路由不能为空字符串", false
	}

	if path[0] != '/' {
		return "web: 路由必须以/开头", false
	}

	if path != "/" {
		if path[len(path)-1] == '/' {
			return "web: 路由不能以/结尾", false
		}

		path = strings.TrimLeft(path, "/")
		pathSegments := strings.Split(path, "/")
		for _, pathSegment := range pathSegments {
			if pathSegment == "" {
				return "web: 路由中不能出现连续的/", false
			}
		}
	}

	return "", true
}

// findRoute 根据给定的HTTP动词和path 在路由树中查找匹配的节点
func (r *router) findRoute(method string, path string) (*matchNode, bool) {
	targetMatchNode := &matchNode{}
	// HTTP动词对应的路由树不存在 直接返回nil false即可
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	target := root
	// 对根节点做特殊处理
	if path == "/" {
		targetMatchNode.node = target
		return targetMatchNode, true
	}

	// 在路由树中逐层查找节点
	path = strings.TrimLeft(path, "/")
	pathSegments := strings.Split(path, "/")
	for _, pathSegment := range pathSegments {
		child, isParam, ok := target.childOf(pathSegment)
		if !ok {
			return nil, false
		}

		if isParam {
			key := strings.TrimPrefix(child.path, ":")
			value := pathSegment
			targetMatchNode.addPathParam(key, value)
		}

		target = child
	}

	targetMatchNode.node = target
	return targetMatchNode, true
}
