package proxy

import (
	"fmt"
	"sync"
)

// ClientRegistry 客户端注册中心
type ClientRegistry struct {
	clients map[string]ThirdPartyClient
	mu      sync.RWMutex
}

// NewClientRegistry 创建客户端注册中心
func NewClientRegistry() *ClientRegistry {
	return &ClientRegistry{
		clients: make(map[string]ThirdPartyClient),
	}
}

// Register 注册客户端
func (r *ClientRegistry) Register(client ThirdPartyClient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := client.GetName()
	if _, exists := r.clients[name]; exists {
		return fmt.Errorf("客户端 %s 已注册", name)
	}

	r.clients[name] = client
	return nil
}

// GetClient 获取客户端
func (r *ClientRegistry) GetClient(name string) (ThirdPartyClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[name]
	if !exists {
		return nil, fmt.Errorf("客户端 %s 不存在", name)
	}

	return client, nil
}

// ListClients 列出所有客户端信息
func (r *ClientRegistry) ListClients() []ClientInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clients := make([]ClientInfo, 0, len(r.clients))
	for _, client := range r.clients {
		clients = append(clients, ClientInfo{
			Name:        client.GetName(),
			RoutePrefix: client.GetRoutePrefix(),
			Description: client.GetDescription(),
			Status:      client.GetStatus(),
		})
	}

	return clients
}

// globalRegistry 全局注册中心
var globalRegistry = NewClientRegistry()

// GetGlobalRegistry 获取全局注册中心
func GetGlobalRegistry() *ClientRegistry {
	return globalRegistry
}



