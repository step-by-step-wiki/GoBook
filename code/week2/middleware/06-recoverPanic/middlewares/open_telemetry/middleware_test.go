package open_telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"web"
)

// Test_MiddlewareBuilder 以创建 web.HandleFunc 的形式,在这个HandleFunc
// 中创建几个具有父子关系的span,并最终在第三方的tracing平台中查看这些span
func Test_MiddlewareBuilder(t *testing.T) {
	// 创建中间件
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := &MiddlewareBuilder{
		Tracer: tracer,
	}
	tracingMiddleware := builder.Build()

	// 创建中间件Option
	middlewareOption := web.ServerWithMiddleware(tracingMiddleware)

	// 创建服务器
	s := web.NewHTTPServer(middlewareOption)

	// 创建HandleFunc
	handleFunc := func(ctx *web.Context) {
		// 创建第1层span
		firstLayerContext1, firstLayerSpan1 := tracer.Start(ctx.Req.Context(), "first_layer_1")
		_, firstLayerSpan2 := tracer.Start(ctx.Req.Context(), "first_layer_2")

		// 创建第2层span 第2层的span是第1层span的子span
		secondLayerContext, secondLayerSpan := tracer.Start(firstLayerContext1, "second_layer")
		// 暂停1s 目的在于: 使得第2层span的开始时间与第1层span的开始时间有较为明显的时间差
		time.Sleep(time.Second)
		// 这里的defer secondLayerSpan.End()其实是可以不写的,因为在第一层span结束的时候,第二层span也会结束

		// 创建第3层的span 第3层的span是第2层span的子span
		_, thirdLayerSpan1 := tracer.Start(secondLayerContext, "third_layer_1")
		// 暂停100ms 先关闭第3层的第1个span
		time.Sleep(100 * time.Millisecond)
		thirdLayerSpan1.End()
		_, thirdLayerSpan2 := tracer.Start(secondLayerContext, "third_layer_2")
		// 暂停300ms 再关闭第3层的第2个span
		// 这样做是为了让第3层的2个span之间是一个有明显时间差的关系
		time.Sleep(300 * time.Millisecond)
		thirdLayerSpan2.End()

		// 再关闭第2层的span
		secondLayerSpan.End()

		// 最后关闭第1层的span
		firstLayerSpan1.End()
		firstLayerSpan2.End()

		ctx.RespJSON(http.StatusOK, User{Name: "test"})
	}

	// 设置TracerProvider
	initZipkin(t)

	// 注册路由并启动服务器
	s.GET("/user", handleFunc)
	s.Start(":8092")
}

type User struct {
	Name string `json:"name"`
}

// initZipkin 初始化zipkin并设置TracerProvider
func initZipkin(t *testing.T) {
	// 要注意这个端口，和 docker-compose 中的保持一致
	exporter, err := zipkin.New(
		"http://localhost:19411/api/v2/spans",
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		t.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
		)),
	)
	otel.SetTracerProvider(tp)
}

// initJeager 初始化jeager并设置TracerProvider
func initJeager(t *testing.T) {
	url := "http://localhost:14268/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		t.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)

	otel.SetTracerProvider(tp)
}
