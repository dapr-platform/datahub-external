package main

import (
	"datahub-external/api"
	_ "datahub-external/docs"
	_ "datahub-external/service"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	PORT         = 80
	BASE_CONTEXT = ""
)

func init() {
	if val := os.Getenv("LISTEN_PORT"); val != "" {
		PORT, _ = strconv.Atoi(val)
	}

	if val := os.Getenv("BASE_CONTEXT"); val != "" {
		BASE_CONTEXT = val
	}
}

// @title 数据底座外部接口代理服务 API
// @version 1.0
// @description 代理第三方数据接口，如绿云酒店管理系统等
// @BasePath /swagger/datahub-external
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API Key认证，也可以使用 Authorization: Bearer <key>
func main() {
	mux := chi.NewRouter()

	// 如果有BASE_CONTEXT，则在该路径下挂载所有路由
	if BASE_CONTEXT != "" {
		mux.Route(BASE_CONTEXT, func(r chi.Router) {
			// 创建子路由器并初始化路由
			subMux := r.(*chi.Mux)
			api.InitRoute(subMux)
			r.Handle("/metrics", promhttp.Handler())
			r.Handle("/swagger*", httpSwagger.WrapHandler)
		})
	} else {
		api.InitRoute(mux)
		mux.Handle("/metrics", promhttp.Handler())
		mux.Handle("/swagger*", httpSwagger.WrapHandler)
	}

	s := daprd.NewServiceWithMux(":"+strconv.Itoa(PORT), mux)
	slog.Info("datahub-external服务启动", "port", PORT)
	if BASE_CONTEXT != "" {
		slog.Info("配置BASE_CONTEXT", "context", BASE_CONTEXT)
		slog.Info("Swagger文档地址", "url", "http://localhost:"+strconv.Itoa(PORT)+BASE_CONTEXT+"/swagger/index.html")
	} else {
		slog.Info("Swagger文档地址", "url", "http://localhost:"+strconv.Itoa(PORT)+"/swagger/index.html")
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		slog.Error("启动失败", "error", err)
		os.Exit(1)
	}
}


