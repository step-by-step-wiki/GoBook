package v6

import (
	"fmt"
	"strings"
)

// router 路由森林 用于支持对路由树的操作
type router struct {
	// trees 路由森林 按HTTP动词组织路由树
	// 该map中 key为HTTP动词 value为路由树的根节点
	// 即: 每个HTTP动词对应一棵路由树 指向每棵路由树的根节点
	trees map[string]*node
}

// newRouter 创建路由森林
func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// AddRoute 注册路由到路由森林中的路由树上
// 其中path为路由的路径.该路径:
// 1. 不得为空字符串
// 2. 必须以"/"开头
// 3. 不能以"/"结尾
// 4. 不能包含连续的"/"
// - 已经注册了的路由,无法被覆盖.例如`/user/home`注册两次,会冲突
// - `06-stringValue`必须以`/`开始并且结尾不能有`/`,中间也不允许有连续的`/`
// - 不能在同一个位置注册不同的参数路由.例如`/user/:id`和`/user/:name`冲突
// - 不能在同一个位置同时注册通配符路由和参数路由.例如`/user/:id`和`/user/*`冲突
// - 同名路径参数,在路由匹配的时候,值会被覆盖.例如`/user/:id/abc/:id`,那么`/user/123/abc/456`,最终`id = 456`
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	// step1. 检测路由是否合规

	// 1.1 检测路由是否为空字符串
	if path == "" {
		panic("web: 路由不能为空字符串")
	}

	// 1.2 检测路由是否以"/"开头
	if path[0] != '/' {
		panic("web: 路由必须以 '/' 开头")
	}

	// 1.3 检测路由是否以"/"结尾
	// Tips: 这个逻辑判断放在根节点的处理后边确实是可以省点代码 但是我认为那样不太好理解
	// Tips: 我认为正常的处理流程是:先判断入参是否合规,再进行后续的逻辑处理.仅当入参合规时,才进行后续的逻辑处理
	// Tips: 因此我把这部分逻辑判断放在根节点的处理前边
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 '/' 结尾")
	}

	// step2. 找到路由树
	root, ok := r.trees[method]
	// 如果没有找到路由树,则创建一棵路由树
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// step3. 判断path是否为根节点 如果是则直接设置HandleFunc并返回即可
	if path == "/" {
		// 判断根节点是否路由冲突
		if root.HandleFunc != nil {
			panic("web: 路由冲突,重复注册路由 [/] ")
		}
		root.HandleFunc = handleFunc
		return
	}

	// step4. 切割path
	// Tips: 去掉前导的"/" 否则直接切割出来的第一个元素为空字符串
	// Tips: 以下代码是老师写的去掉前导的"/"的方式 我认为表达力有点弱 但是性能应该会好于strings.TrimLeft
	// Tips: 以下代码会有问题,因为假如前导字符不是"/" 则不该被去掉
	// path = path[1:]
	path = strings.TrimLeft(path, "/")
	segments := strings.Split(path, "/")

	// step3. 为路由树添加路由
	// Tips: 此处我认为用target指代要添加路由的节点更好理解
	target := root
	for _, segment := range segments {
		// 若切割后的路由段为空字符串,则说明路由中有连续的"/"
		if segment == "" {
			panic("web: 路由中不得包含连续的'/'")
		}

		// 如果路由树中途有节点没有创建,则创建该节点;
		// 如果路由树中途存在子节点,则找到该子节点
		child := target.childOrCreate(segment)
		// 继续为子节点创建子节点
		target = child
	}

	// 判断普通节点是否路由冲突
	if target.HandleFunc != nil {
		panic(fmt.Sprintf("web: 路由冲突,重复注册路由 [%s] ", path))
	}

	// 为目标节点设置HandleFunc
	target.HandleFunc = handleFunc
}

// findRoute 根据给定的HTTP方法和路由路径,在路由森林中查找对应的节点
// 若该节点为参数路径节点,则不仅返回该节点,还返回参数名和参数值
// 否则,仅返回该节点
func (r *router) findRoute(method string, path string) (*matchNode, bool) {
	targetMatchNode := &matchNode{}
	root, ok := r.trees[method]
	// 给定的HTTP动词在路由森林中不存在对应的路由树,则直接返回false
	if !ok {
		return nil, false
	}

	// 对根节点做特殊处理
	if path == "/" {
		targetMatchNode.node = root
		return targetMatchNode, true
	}

	// 给定的HTTP动词在路由森林中存在对应的路由树,则在该路由树中查找对应的节点
	// 去掉前导和后置的"/"
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	// Tips: 同样的 这里我认为用target作为变量名表现力更强
	target := root

	for _, segment := range segments {
		child, isParamChild, found := target.childOf(segment)
		// 如果在当前节点的子节点映射中没有找到对应的子节点,则直接返回
		if !found {
			return nil, false
		}

		// 若当前节点为参数节点,则将参数名和参数值保存到targetMatchNode中
		if isParamChild {
			// 参数名是形如 :id 的格式, 因此需要去掉前导的:
			name := child.path[1:]
			// 参数值就是当前路由路径中的路由段
			value := segment
			targetMatchNode.addPathParams(name, value)
		}

		// 如果在当前节点的子节点映射中找到了对应的子节点,则继续在该子节点中查找
		target = child
	}

	// 如果找到了对应的节点,则返回该节点
	// Tips: 此处有2种设计 一种是用标量表示是否找到了子节点
	// Tips: 另一种是 return target, target.HandleFunc != nil
	// Tips: 这种返回就表示找到了子节点且子节点必然有对应的业务处理函数
	// 此处我倾向用第1种设计 因为方法名叫findRoute,表示是否找到节点的意思.而非表示是否找到了一个有对应的业务处理函数的节点
	targetMatchNode.node = target
	return targetMatchNode, true
}
