package bootstrap

import (
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	apispec "github.com/konlyk/go_api_skeleton/api"
	"github.com/konlyk/go_api_skeleton/api/openapi"
	"github.com/konlyk/go_api_skeleton/controller/middleware"
)

func SetupRouter(app *Application) (*gin.Engine, error) {
	engine := gin.New()

	metrics := middleware.NewAPIMetrics(app.Observability.MetricsRegistry)
	engine.Use(
		middleware.Recoverer(app.Observability.Logger),
		middleware.RequestLogger(app.Observability.Logger),
		metrics.Handler(),
		otelgin.Middleware(app.Config.ServiceName),
	)

	engine.GET("/metrics", gin.WrapH(promhttp.HandlerFor(app.Observability.MetricsRegistry, promhttp.HandlerOpts{})))

	swagger, err := apispec.LoadBundledSpec()
	if err != nil {
		return nil, err
	}
	swagger.Servers = nil

	apiGroup := engine.Group("/")
	apiGroup.Use(
		ginmiddleware.OapiRequestValidatorWithOptions(swagger, &ginmiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: newOpenAPIAuthenticationFunc(app.Config.PrivateAPIToken),
			},
		}),
	)

	controllers := app.Controllers
	strictServer := openapi.NewStrictHandler(
		controllers,
		nil,
	)

	openapi.RegisterHandlersWithOptions(apiGroup, strictServer, openapi.GinServerOptions{})

	return engine, nil
}
