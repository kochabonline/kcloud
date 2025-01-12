package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kochabonline/kit/transport/http"
	"github.com/kochabonline/kit/transport/http/middleware"

	"github.com/kochabonline/kcloud/apps/system/auth/jwt"
	"github.com/kochabonline/kcloud/config"
	"github.com/kochabonline/kcloud/docs"
)

// Handler 路由注册接口
type Handler interface {
	Register(r gin.IRouter)
}

func NewHttpServer(config *config.Config, jwtController jwt.Interface, handlers ...any) *http.Server {
	// 跳过需要验证header的路径, 例如登录接口
	skippedPathPrefixes := []string{
		"/health",
		"/swagger",
		"/metrics",
		"/api/v1/auth/login",
		"/api/v1/notifier/message/send",
	}

	gin.SetMode(config.Http.Level)
	r := gin.New()
	r.Use(
		middleware.GinLogger(),
		middleware.GinRecovery(),
		middleware.AuthWithConfig(middleware.AuthConfig{
			SkippedPathPrefixes: skippedPathPrefixes,
			Validate:            jwtController.Validate,
		}),
	)
	apiV1 := r.Group("/api/v1")
	// 注册所有路由
	for _, handler := range handlers {
		if h, ok := handler.(Handler); ok {
			h.Register(apiV1)
		}
	}

	// swagger文档
	docs.SwaggerInfo.Title = "kcloud API"
	docs.SwaggerInfo.Description = "This is a sample server kcloud server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	server := http.NewServer(config.Http.Addr(), r,
		http.WithSwagOptions(http.SwagOption{Enabled: config.Http.Swagger.Enabled}),
		http.WithHealthOptions(http.HealthOption{Enabled: config.Http.Health.Enabled}),
		http.WithMetricsOptions(http.MetricsOption{Enabled: config.Http.Metrics.Enabled}),
	)

	return server
}
