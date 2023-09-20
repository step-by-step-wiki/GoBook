package staticMatchingTestCase

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

// TestNode 测试路由树节点
// 由于此处我们要测试的是路由树的结构,因此不需要在测试路由树节点中添加路由处理函数
// 调用AddRoute时写死一个HandleFunc即可
type TestNode struct {
	method string
	path   string
}

// TestRouter_AddRoute 测试路由注册功能的结果是否符合预期
func TestRouter_AddRoute(t *testing.T) {
	// step1. 构造路由树
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
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	r := newRouter()
	mockHandleFunc := func(ctx Context) {}

	for _, testRoute := range testRoutes {
		r.AddRoute(testRoute.method, testRoute.path, mockHandleFunc)
	}

	// step2. 验证路由树 断言二者是否相等
	wantRouter := &router{
		trees: map[string]*node{
			// GET方法路由树
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": &node{
								path:     "home",
								children: nil,
								// 注意路由是/user/home 因此只有最深层的节点才有handleFunc
								// /user和/ 都是没有handleFunc的
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

			// POST方法路由树
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
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
					"login": &node{
						path:       "login",
						children:   nil,
						HandleFunc: mockHandleFunc,
					},
				},
				HandleFunc: nil,
			},
		},
	}

	// HandleFunc类型是方法,方法不可比较,因此只能比较两个路由树的结构是否相等
	// assert.Equal(t, wantRouter, r)

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)
}

// equal 比较两个路由森林是否相等
// msg: 两个路由森林不相等时的错误信息
// ok: 两个路由森林是否相等
func (r *router) equal(y *router) (msg string, ok bool) {
	// 如果目标路由森林为nil 则不相等
	if y == nil {
		return fmt.Sprintf("目标路由森林为nil"), false
	}

	// 如果两个路由森林中的路由树数量不同 则不相等
	rTreesNum := len(r.trees)
	yTreesNum := len(y.trees)
	if rTreesNum != yTreesNum {
		return fmt.Sprintf("路由森林中的路由树数量不相等,源路由森林有 %d 棵路由树, 目标路由森林有 %d 棵路由树", rTreesNum, yTreesNum), false
	}

	for method, tree := range r.trees {
		dstTree, ok := y.trees[method]

		// 如果目标router中没有对应HTTP方法的路由树 则不相等
		if !ok {
			return fmt.Sprintf("目标 router 中没有HTTP方法 %s的路由树", method), false
		}

		// 比对两棵路由树的结构是否相等
		msg, equal := tree.equal(dstTree)
		if !equal {
			return method + "-" + msg, false
		}
	}
	return "", true
}

// equal 比较两棵路由树是否相等
// msg: 两棵路由树不相等时的错误信息
// ok: 两棵路由树是否相等
func (n *node) equal(y *node) (msg string, ok bool) {
	// 如果目标节点为nil 则不相等
	if y == nil {
		return fmt.Sprintf("目标节点为nil"), false
	}

	// 如果两个节点的path不相等 则不相等
	if n.path != y.path {
		return fmt.Sprintf("两个节点的path不相等,源节点的path为 %s,目标节点的path为 %s", n.path, y.path), false
	}

	// 若两个节点的子节点数量不相等 则不相等
	nChildrenNum := len(n.children)
	yChildrenNum := len(y.children)
	if nChildrenNum != yChildrenNum {
		return fmt.Sprintf("两个节点的子节点数量不相等,源节点的子节点数量为 %d,目标节点的子节点数量为 %d", nChildrenNum, yChildrenNum), false
	}

	// 若两个节点的handleFunc类型不同 则不相等
	nHandler := reflect.ValueOf(n.HandleFunc)
	yHandler := reflect.ValueOf(y.HandleFunc)
	if nHandler != yHandler {
		return fmt.Sprintf("%s节点的handleFunc不相等,源节点的handleFunc为 %v,目标节点的handleFunc为 %v", n.path, nHandler.Type().String(), yHandler.Type().String()), false
	}

	// 比对两个节点的子节点映射是否相等
	for path, child := range n.children {
		dstChild, ok := y.children[path]
		// 如果源节点的子节点中 存在目标节点没有的子节点 则不相等
		if !ok {
			return fmt.Sprintf("目标节点的子节点中没有path为 %s 的子节点", path), false
		}

		// 比对两个子节点是否相等
		msg, equal := child.equal(dstChild)
		if !equal {
			return msg, false
		}
	}

	return "", true
}

func TestRouter_AddRoute2(t *testing.T) {}
