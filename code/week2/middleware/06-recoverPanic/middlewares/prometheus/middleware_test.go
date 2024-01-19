package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
	"web"
)

// Test_MiddlewareBuilder 测试MiddlewareBuilder
func Test_MiddlewareBuilder(t *testing.T) {
	// 创建中间件构建器
	builder := &MiddlewareBuilder{
		// Tips: 此处的Namespace、Subsystem、Name、Help的值均不可以出现-
		Namespace: "my_framework",
		Subsystem: "web",
		Name:      "http_response",
		Help:      "metric_help",
	}

	// 创建中间件Option
	options := web.ServerWithMiddleware(builder.Build())

	// 创建服务器
	server := web.NewHTTPServer(options)

	// 创建HandleFunc
	handleFunc := func(ctx *web.Context) {
		// 随机sleep 1-1001ms 模拟请求时长
		randomValue := rand.Intn(1000) + 1
		time.Sleep(time.Duration(randomValue) * time.Millisecond)

		ctx.RespJSON(http.StatusAccepted, &User{Name: "Tom"})
	}

	// 监听另外一个端口 该端口用于暴露指标让prometheus采集
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 通常用于采集指标的端口不会和业务端口混合在一起
		// 因为业务端口通常是需要暴露到外网的 但是指标采集端口通常是不需要的
		http.ListenAndServe(":8082", nil)
	}()

	// 注册路由并启动服务器
	server.GET("/user", handleFunc)
	server.Start(":8080")
}

type User struct {
	Name string
}
