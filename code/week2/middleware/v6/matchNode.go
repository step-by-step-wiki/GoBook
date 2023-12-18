package v6

// matchNode 用于保存匹配到的节点与路径参数
type matchNode struct {
	// node 匹配到的节点
	node *node
	// pathParams 路径参数 若该节点不是参数节点,则该字段值为nil
	pathParams map[string]string
}

// addPathParams 用于添加路径参数
func (m *matchNode) addPathParams(name string, value string) {
	// Tips: 这里作为框架的设计者 你是没法确定用户会注册多少个参数路由的
	// Tips: 因此给不给容量意义不大
	if m.pathParams == nil {
		m.pathParams = map[string]string{}
	}
	m.pathParams[name] = value
}
