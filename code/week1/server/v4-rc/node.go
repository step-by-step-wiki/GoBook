package v4_rc

import (
	"fmt"
	"strings"
)

// node 路由树中的节点
type node struct {
	path          string           // path 路由路径
	children      map[string]*node // children 子节点 key为子节点的路由路径 value为路径对应子节点
	wildcardChild *node            // wildcardChild 通配符子节点
	paramChild    *node            // paramChild 参数路由子节点
	HandleFunc                     // HandleFunc 路由对应的处理函数
}

// findOrCreate 本方法用于根据给定的path值 在当前节点的子节点中查找path为给定path值的节点
// 找到则返回 未找到则创建
func (n *node) findOrCreate(segment string) *node {
	// 若路径以:开头 则查找或创建参数子节点
	if strings.HasPrefix(segment, ":") {
		if n.wildcardChild != nil {
			msg := fmt.Sprintf("web: 非法路由,节点 %s 已有通配符路由.不允许同时注册通配符路由和参数路由", n.path)
			panic(msg)
		}

		if n.paramChild == nil {
			n.paramChild = &node{
				path: segment,
			}
		}
		return n.paramChild
	}

	// 若路径为* 则查找或创建通配符子节点
	if segment == "*" {
		if n.paramChild != nil {
			msg := fmt.Sprintf("web: 非法路由,节点 %s 已有参数路由.不允许同时注册通配符路由和参数路由", n.path)
			panic(msg)
		}

		if n.wildcardChild == nil {
			n.wildcardChild = &node{
				path: "*",
			}
		}
		return n.wildcardChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}

	target, exist := n.children[segment]
	if !exist {
		// 当前节点的子节点映射中不存在目标子节点 则创建目标子节点 将子节点加入当前节点的子节点映射后返回
		target = &node{
			path: segment,
		}
		n.children[segment] = target
		return target
	}

	// 当前节点的子节点映射中存在目标子节点 则直接返回
	return target
}

// childOf 本方法用于根据给定的path值 在当前节点的子节点映射中查找对应的子节点
// 若未在当前节点的子节点映射中查找到path对应的节点 则尝试查找当前节点的参数子节点
// 若未查找到当前节点的参数子节点 则尝试查找当前节点的通配符子节点
func (n *node) childOf(path string) (targetNode *node, isParamNode bool, isFound bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		return n.wildcardChild, false, n.wildcardChild != nil
	}

	child, found := n.children[path]
	if !found {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		return n.wildcardChild, false, n.wildcardChild != nil
	}

	return child, false, true
}
