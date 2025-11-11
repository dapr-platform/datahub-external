package proxy

import "context"

// ThirdPartyClient 第三方接口客户端接口
type ThirdPartyClient interface {
	// GetName 获取客户端名称
	GetName() string

	// GetRoutePrefix 获取支持的路由前缀
	GetRoutePrefix() string

	// GetDescription 获取客户端描述
	GetDescription() string

	// HandleRequest 处理请求
	// path: 请求路径(不包含前缀)
	// params: 查询参数
	HandleRequest(ctx context.Context, path string, params map[string]string) (interface{}, error)

	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error

	// GetStatus 获取客户端状态
	GetStatus() string
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Name        string `json:"name"`
	RoutePrefix string `json:"route_prefix"`
	Description string `json:"description"`
	Status      string `json:"status"`
}














