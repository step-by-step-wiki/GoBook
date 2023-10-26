package v4_rc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

type TestNode struct {
	method string
	path   string
}

func TestRouter_AddRoute(t *testing.T) {
	// step1. 构造路由森林
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
	}
	targetRouter := newRouter()
	mockHandleFunc := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		targetRouter.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 验证路由树
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"home": &node{
								path:       "home",
								children:   nil,
								HandleFunc: mockHandleFunc,
							},
						},
						HandleFunc: mockHandleFunc,
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:       "detail",
								children:   nil,
								HandleFunc: mockHandleFunc,
							},
						},
						HandleFunc: nil,
					},
				},
				HandleFunc: mockHandleFunc,
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"login": &node{
						path:       "login",
						children:   nil,
						HandleFunc: mockHandleFunc,
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:       "create",
								children:   nil,
								HandleFunc: mockHandleFunc,
							},
						},
						HandleFunc: nil,
					},
				},
				HandleFunc: nil,
			},
		},
	}

	msg, ok := wantRouter.equal(targetRouter)
	assert.True(t, ok, msg)
}

func (r *router) equal(target *router) (msg string, ok bool) {
	// step1. 比较路由森林中的路由树数量
	wantLen := len(r.trees)
	targetLen := len(target.trees)

	if wantLen != targetLen {
		msg = fmt.Sprintf("路由森林中的路由树数量不等, 期望路由树的数量为: %d, 目标路由树的数量为: %d", wantLen, targetLen)
		return msg, false
	}

	// step2. 比对2个路由森林中的路由树HTTP动词是否相同
	for method, tree := range r.trees {
		dstTree, ok := target.trees[method]
		if !ok {
			msg = fmt.Sprintf("目标路由森林中不存在HTTP动词为: %s 的路由树", method)
			return msg, false
		}

		// step3. 比对2个路由树中的结构是否相同
		msg, ok = tree.equal(dstTree)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func (n *node) equal(target *node) (msg string, ok bool) {
	// step1. 目标节点为空 则必然不等
	if target == nil {
		msg = "目标节点为nil"
		return msg, false
	}

	// step2. 比较节点的路径是否相同
	if n.path != target.path {
		msg = fmt.Sprintf("节点的路径不等, 期望节点的路径为: %s, 目标节点的路径为: %s", n.path, target.path)
		return msg, false
	}

	// step3. 比较2个节点的子节点数量是否相同
	if len(n.children) != len(target.children) {
		msg = fmt.Sprintf("节点的子节点数量不等, 期望节点的子节点数量为: %d, 目标节点的子节点数量为: %d", len(n.children), len(target.children))
		return msg, false
	}

	// step4. 比对2个节点的参数子节点是否相同
	if n.paramChild != nil {
		if target.paramChild == nil {
			msg = fmt.Sprintf("目标节点的参数节点为空")
			return msg, false
		}
		_, equal := n.paramChild.equal(target.paramChild)
		if !equal {
			msg = fmt.Sprintf("期望节点 %s 的参数子节点与目标节点 %s 的参数子节点不等", n.path, target.path)
			return msg, false
		}
	}

	// step5. 比较2个节点的通配符子节点是否相同
	if n.wildcardChild != nil {
		if target.wildcardChild == nil {
			msg = fmt.Sprintf("目标节点的通配符子节点为空")
			return msg, false
		}
		_, equal := n.wildcardChild.equal(target.wildcardChild)
		if !equal {
			msg = fmt.Sprintf("期望节点 %s 的通配符子节点与目标节点 %s 的通配符子节点不等", n.path, target.path)
			return msg, false
		}
	}

	// step5. 比较2个节点的处理函数是否相同
	wantHandler := reflect.ValueOf(n.HandleFunc)
	targetHandler := reflect.ValueOf(target.HandleFunc)
	if wantHandler != targetHandler {
		msg = fmt.Sprintf("节点的处理函数不等, 期望节点 %s 的处理函数为: %v, 目标节点 %s 的处理函数为: %v", n.path, wantHandler, target.path, targetHandler)
		return msg, false
	}

	// step6. 比较2个节点的子节点是否相同
	for path, child := range n.children {
		// step6.1 比对2个节点的子节点的路径是否相同
		dstChild, exist := target.children[path]
		if !exist {
			msg = fmt.Sprintf("目标节点中不存在路径为: %s 的子节点", path)
			return msg, false
		}

		// step6.2 对路径相同的子节点递归比对
		msg, equal := child.equal(dstChild)
		if !equal {
			return msg, false
		}
	}

	return "", true
}

func TestRouter_Illegal_Path(t *testing.T) {
	r := newRouter()
	mockHandle := func(ctx *Context) {}

	// 路由为空的测试用例
	nilPathFunc := func() {
		r.addRoute(http.MethodGet, "", mockHandle)
	}
	assert.Panicsf(t, nilPathFunc, "web: 路由不能为空字符串")

	// 路由不是以`/`开头的测试用例
	incorrectFirstCharacter := func() {
		r.addRoute(http.MethodGet, "login", mockHandle)
	}
	assert.Panicsf(t, incorrectFirstCharacter, "web: 路由必须以/开头")

	// 路由以`/`结尾的测试用例
	incorrectLastCharacter := func() {
		r.addRoute(http.MethodGet, "/login/", mockHandle)
	}
	assert.Panicsf(t, incorrectLastCharacter, "web: 路由不能以/结尾")

	// 路由中出现了连续`/`的测试用例
	continuousSeparator := func() {
		r.addRoute(http.MethodGet, "/a//b", mockHandle)
	}
	assert.Panicsf(t, continuousSeparator, "web: 路由不能出现多个连续的/")

	// 路由重复注册的测试用例
	// 根节点路由重复注册
	r.addRoute(http.MethodGet, "/", mockHandle)
	repeatRegisterRoute := func() {
		r.addRoute(http.MethodGet, "/", mockHandle)
	}
	assert.Panicsf(t, repeatRegisterRoute, "web: 路由冲突,重复注册路由 [/]")

	// 普通节点路由重复注册
	r.addRoute(http.MethodGet, "/user/login", mockHandle)
	repeatRegisterRoute = func() {
		r.addRoute(http.MethodGet, "/user/login", mockHandle)
	}
	assert.Panicsf(t, repeatRegisterRoute, "web: 路由冲突,重复注册路由 [/user/login]")
}

func TestRouter_FindRoute(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
	}
	r := newRouter()
	mockHandle := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandle)
	}

	// step2. 构造测试用例
	testCases := []struct {
		// name 子测试用例的名称
		name string
		// method HTTP动词
		method string
		// path 路由路径
		path string
		// isFound 是否找到了节点
		isFound bool
		// wantNode 期望的路由节点
		wantNode *matchNode
	}{
		{
			name:     "Method not found",
			method:   http.MethodDelete,
			path:     "/user",
			isFound:  false,
			wantNode: nil,
		},
		{
			name:    "completely match",
			method:  http.MethodGet,
			path:    "/order/detail",
			isFound: true,
			wantNode: &matchNode{
				node: &node{
					path:       "detail",
					children:   nil,
					HandleFunc: mockHandle,
				},
				pathParams: nil,
			},
		},
		{
			name:    "nil handle func",
			method:  http.MethodGet,
			path:    "/order",
			isFound: true,
			wantNode: &matchNode{
				node: &node{
					path: "order",
					children: map[string]*node{
						"detail": &node{
							path:       "detail",
							children:   nil,
							HandleFunc: mockHandle,
						},
					},
					HandleFunc: nil,
				},
				pathParams: nil,
			},
		},
		{
			name:    "root node",
			method:  http.MethodGet,
			path:    "/",
			isFound: true,
			wantNode: &matchNode{
				node: &node{
					path: "/",
					children: map[string]*node{
						"user": &node{
							path:       "user",
							children:   nil,
							HandleFunc: mockHandle,
						},
						"order": &node{
							path: "order",
							children: map[string]*node{
								"detail": &node{
									path:       "detail",
									children:   nil,
									HandleFunc: mockHandle,
								},
							},
							HandleFunc: nil,
						},
					},
					HandleFunc: mockHandle,
				},
				pathParams: nil,
			},
		},
		{
			name:     "path not found",
			method:   http.MethodGet,
			path:     "/login",
			isFound:  false,
			wantNode: nil,
		},
	}

	// step3. 测试是否找到节点
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundNode, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.isFound, found)
			// 3.1 判断在路由树中是否找到了节点
			if !found {
				return
			}

			// 3.2 判断找到的节点和预期的节点是否相同
			msg, equal := testCase.wantNode.node.equal(foundNode.node)
			assert.True(t, equal, msg)
		})
	}
}

type TestCaseNode struct {
	name     string // name 子测试用例的名称
	method   string // method HTTP动词
	path     string // path 路由路径
	isFound  bool   // isFound 是否找到了节点
	wantNode *node  // wantNode 期望的路由节点
}

func TestRouter_wildcard(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		// 普通节点的通配符子节点
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		// 根节点的通配符子节点
		{
			method: http.MethodGet,
			path:   "/*",
		},
		// 通配符子节点的通配符子节点
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		// 通配符子节点的普通子节点
		{
			method: http.MethodGet,
			path:   "/*/get",
		},
		// 通配符子节点的普通子节点的通配符子节点
		{
			method: http.MethodGet,
			path:   "/*/order/*",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 断言路由树
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path:     "order",
						children: nil,
						wildcardChild: &node{
							path:          "*",
							children:      nil,
							wildcardChild: nil,
							HandleFunc:    mockHandleFunc,
						},
						HandleFunc: nil,
					},
				},
				wildcardChild: &node{
					path: "*",
					children: map[string]*node{
						"get": &node{
							path:          "get",
							children:      nil,
							wildcardChild: nil,
							HandleFunc:    mockHandleFunc,
						},
						"order": &node{
							path:     "order",
							children: nil,
							wildcardChild: &node{
								path:          "*",
								children:      nil,
								wildcardChild: nil,
								HandleFunc:    mockHandleFunc,
							},
							HandleFunc: nil,
						},
					},
					wildcardChild: &node{
						path:          "*",
						children:      nil,
						wildcardChild: nil,
						HandleFunc:    mockHandleFunc,
					},
					HandleFunc: mockHandleFunc,
				},
				HandleFunc: nil,
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)
}

func TestRouter_FindRoute_Wildcard(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/order/create",
		},
	}

	r := newRouter()
	mockHandle := func(ctx *Context) {}

	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandle)
	}

	// step2. 构造测试用例
	testCases := []struct {
		name     string
		method   string
		path     string
		isFound  bool
		wantNode *matchNode
	}{
		{
			name:    "普通节点的通配符子节点",
			method:  http.MethodGet,
			path:    "/order/detail",
			isFound: true,
			wantNode: &matchNode{
				node: &node{
					path:          "*",
					children:      nil,
					wildcardChild: nil,
					HandleFunc:    mockHandle,
				},
				pathParams: nil,
			},
		},
		{
			name:    "普通节点下通配符子节点和普通子节点共存",
			method:  http.MethodGet,
			path:    "/order/create",
			isFound: true,
			wantNode: &matchNode{
				node: &node{
					path:          "create",
					children:      nil,
					wildcardChild: nil,
					HandleFunc:    mockHandle,
				},
				pathParams: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			targetNode, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.isFound, found)
			if !found {
				return
			}

			msg, equal := testCase.wantNode.node.equal(targetNode.node)
			assert.True(t, equal, msg)
		})
	}
}

func TestRouter_addParamRoute(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []TestNode{
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 验证路由树
	wantRoute := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:          "detail",
								children:      nil,
								wildcardChild: nil,
								paramChild: &node{
									path:          ":id",
									children:      nil,
									wildcardChild: nil,
									paramChild:    nil,
									HandleFunc:    mockHandleFunc,
								},
								HandleFunc: nil,
							},
						},
						wildcardChild: nil,
						paramChild:    nil,
						HandleFunc:    nil,
					},
				},
				wildcardChild: nil,
				paramChild:    nil,
				HandleFunc:    nil,
			},
		},
	}

	msg, equal := wantRoute.equal(r)
	assert.True(t, equal, msg)
}

func TestRouter_findRoute_param_and_wildcard_coexist(t *testing.T) {
	// step1. 注册通配符路由
	r := newRouter()
	mockHandleFunc := func(ctx *Context) {}
	r.addRoute(http.MethodGet, "/user/*", mockHandleFunc)

	// step2. 断言非法注册
	panicFunc := func() {
		r.addRoute(http.MethodGet, "/user/:id", mockHandleFunc)
	}

	assert.Panicsf(t, panicFunc, "web: 非法路由,节点 detail 已有通配符路由.不允许同时注册通配符路由和参数路由")
}

func TestRouter_findParamRoute(t *testing.T) {
	// step1. 构造路由树
	testRoutes := []*TestNode{
		{
			method: http.MethodGet,
			path:   "/order/:id",
		},
		{
			method: http.MethodGet,
			path:   "/user/:id/detail",
		},
	}

	r := newRouter()
	mockHandle := func(ctx *Context) {}
	for _, testRoute := range testRoutes {
		r.addRoute(testRoute.method, testRoute.path, mockHandle)
	}

	// step2. 构造测试用例
	testCases := []struct {
		name     string
		method   string
		path     string
		isFound  bool
		wantNode *matchNode
	}{
		{
			name:    "param route",
			method:  http.MethodGet,
			path:    "/order/5",
			isFound: true,
			wantNode: &matchNode{
				node: &node{
					path:          ":id",
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					HandleFunc:    mockHandle,
				},
				pathParams: map[string]string{
					"id": "5",
				},
			},
		},
		{
			name:    "param route",
			method:  http.MethodGet,
			path:    "/user/1/detail",
			isFound: true,
			wantNode: &matchNode{
				node: &node{
					path:          "detail",
					children:      nil,
					wildcardChild: nil,
					paramChild:    nil,
					HandleFunc:    mockHandle,
				},
				pathParams: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			findNode, found := r.findRoute(testCase.method, testCase.path)
			assert.True(t, found, "节点未找到")
			if !found {
				return
			}
			msg, equal := testCase.wantNode.node.equal(findNode.node)
			assert.True(t, equal, msg)
		})
	}
}
