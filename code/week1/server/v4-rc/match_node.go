package v4_rc

type matchNode struct {
	node       *node             // node 命中的节点
	pathParams map[string]string // pathParams 节点对应的路由参数 其中key为参数名 value为参数值
}

// addPathParam 用于添加路径参数
func (m *matchNode) addPathParam(key string, value string) {
	if m.pathParams == nil {
		m.pathParams = make(map[string]string)
	}
	m.pathParams[key] = value
}
