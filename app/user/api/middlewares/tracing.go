package middlewares

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/api/global"
)

func Tracing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1, // 全部采样
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d",
					global.RemoteConfig.JaegerInfo.Host,
					global.RemoteConfig.JaegerInfo.Port,
				),
			},
			ServiceName: global.RemoteConfig.ServiceInfo.Name,
		}
		tracer, c, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			global.Logger.Fatal("创建 tracer 失败", zap.Error(err))
		}
		defer func(c io.Closer) {
			_ = c.Close()
		}(c)
		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()
		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", startSpan)
		ctx.Next()
	}
}
