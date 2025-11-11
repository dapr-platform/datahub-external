package api

import (
	"datahub-external/api/controllers"
	"datahub-external/api/middleware"
	"datahub-external/service/proxy"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// InitRoute 初始化所有API路由
func InitRoute(r *chi.Mux) {
	// 基础中间件
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS配置
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API Key认证中间件
	apiKeyAuth := middleware.NewApiKeyAuthMiddleware()
	r.Use(apiKeyAuth.Middleware)

	// 健康检查(白名单,无需认证)
	healthController := controllers.NewHealthController()
	r.Get("/health", healthController.Health)
	r.Get("/ready", healthController.Ready)

	// 数据源列表接口
	proxyController := controllers.NewProxyController(proxy.GetGlobalRegistry())
	r.Get("/datasources", proxyController.ListDataSources)

	// 绿云接口(需要认证)
	lvyunController := controllers.NewLvyunController(proxy.GetGlobalRegistry())
	r.Route("/lvyun", func(r chi.Router) {
		r.Get("/health", lvyunController.HealthCheck)
		r.Get("/info", lvyunController.GetInfo)
		r.Get("/reservations", lvyunController.GetReservations)
		r.Get("/registrations", lvyunController.GetRegistrations)
		r.Get("/checkouts", lvyunController.GetCheckouts)
		r.Get("/business-report", lvyunController.GetBusinessReport)
	})

	// PS系统接口(需要认证)
	psController := controllers.NewPSController(proxy.GetGlobalRegistry())
	r.Route("/ps", func(r chi.Router) {
		r.Get("/health", psController.HealthCheck)
		r.Get("/info", psController.GetInfo)
		r.Get("/family-members", psController.GetFamilyMembers)
		r.Post("/trigger-sync", psController.TriggerSync)
	})

	// 404处理
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, map[string]interface{}{
			"status":  http.StatusNotFound,
			"message": "接口不存在",
			"error":   "Not Found",
		})
	})
}
