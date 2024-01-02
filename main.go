// Code generated by hertz generator.

package main

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/provider"
	"github.com/CloudStriver/go-pkg/utils/hertz/middleware"
	"github.com/CloudStriver/go-pkg/utils/util/log"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func Init() {
	provider.Init()
	hlog.SetLogger(log.NewHlogLogger())
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(b3.New(), propagation.Baggage{}, propagation.TraceContext{}))
}

func main() {
	Init()
	c := provider.Get().Config
	tracer, cfg := tracing.NewServerTracer()
	h := server.New(
		server.WithHostPorts(c.ListenOn),
		server.WithTracer(prometheus.NewServerTracer(":9091", "/server/metrics")),
		tracer,
	)
	h.Use(tracing.ServerMiddleware(cfg), middleware.EnvironmentMiddleware, recovery.Recovery(), func(ctx context.Context, c *app.RequestContext) {
		ctx = adaptor.InjectContext(ctx, c)
		c.Next(ctx)
	})

	register(h)
	log.Info("server start")
	h.Spin()
}
