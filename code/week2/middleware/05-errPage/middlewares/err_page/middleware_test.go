package err_page

import (
	"testing"
	"web"
)

// Test_MiddlewareBuilder 测试错误页面中间件构造器
func Test_MiddlewareBuilder(t *testing.T) {
	// 创建中间件构建器
	builder := NewMiddlewareBuilder()
	builder.
		AddCode(404, []byte(`
<html>
	<head>
		<title>404 Not Found</title>
	</head>

	<body>
		<h1>404 Not Found</h1>
	</body>
</html>
`)).
		AddCode(500, []byte(`
<html>
	<head>
		<title>500 Internal Server Error</title>
	</head>

	<body>
		<h1>500 Internal Server Error</h1>
	</body>
</html>
`))

	// 创建中间件Option
	options := web.ServerWithMiddleware(builder.Build())

	// 创建服务器
	server := web.NewHTTPServer(options)

	// 启动服务器
	server.Start(":8080")
}
