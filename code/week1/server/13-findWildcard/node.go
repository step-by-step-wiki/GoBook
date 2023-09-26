package findWildcard

// node 路由树的节点
type node struct {
	// path 当前节点的路径
	path string

	// children 子路由路径到子节点的映射
	children map[string]*node

	// wildcardChild 通配符子节点
	wildcardChild *node

	// HandleFunc 路由对应的业务逻辑
	HandleFunc
}

// childOrCreate 本方法用于在节点上获取给定的子节点,如果给定的子节点不存在则创建
func (n *node) childOrCreate(segment string) *node {
	// 若路径为通配符 则查找当前节点的通配符子节点 或创建一个当前节点的通配符子节点 并返回
	if segment == "*" {
		if n.wildcardChild == nil {
			n.wildcardChild = &node{
				path: segment,
			}
		}
		return n.wildcardChild
	}

	// 如果当前节点的子节点映射为空 则创建一个子节点映射
	if n.children == nil {
		n.children = map[string]*node{}
	}

	res, ok := n.children[segment]
	// 如果没有找到子节点,则创建一个子节点
	// 否则返回找到的子节点
	if !ok {
		res = &node{
			path: segment,
		}
		n.children[segment] = res
	}
	return res
}

// childOf 根据给定的path在当前节点的子节点映射中查找对应的子节点(即匹配到了静态路由)
// 若未在子节点映射中找到对应子节点 则尝试返回当前节点的通配符子节点
func (n *node) childOf(path string) (child *node, found bool) {
	// 当前节点的子节点映射为空 则有可能匹配到通配符节点
	if n.children == nil {
		return n.wildcardChild, n.wildcardChild != nil
	}

	// 在子当前节点的节点映射中查找对应的子节点 若未找到同样尝试返回当前节点的通配符子节点
	child, found = n.children[path]
	if !found {
		return n.wildcardChild, n.wildcardChild != nil
	}

	// 找到了对应的子节点 则返回该子节点
	return child, found
}
