package initialize

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/jimyag/shop/app/order/rpc/global"
)

func InitJaeger() (opentracing.Tracer, io.Closer, error) {
	// 初始化jaeger
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
	tracer, cl, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	return tracer, cl, err
}
