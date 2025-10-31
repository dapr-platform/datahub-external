# datahub-external

第三方接口代理服务，为 `datahub-service` 提供外部接口访问能力。

## 项目简介

`datahub-external` 是一个轻量级的Go微服务，专门用于代理和简化对复杂第三方接口的访问。通过统一的API接口，可以方便地接入各种外部系统，如绿云酒店管理系统等。

## 架构特点

- **轻量级设计**: 无数据库依赖，纯内存管理
- **可插拔架构**: 基于接口设计，易于扩展新的第三方客户端
- **统一鉴权**: 使用API Key进行统一认证
- **自动管理**: 自动处理Session刷新等复杂逻辑
- **文档完善**: 集成Swagger API文档

## 目录结构

```
datahub-external/
├── main.go                 # 程序入口
├── go.mod                  # Go模块依赖
├── Dockerfile              # 容器化构建
├── api/                    # API层
│   ├── routes.go           # 路由定义
│   ├── controllers/        # 控制器
│   │   ├── health_controller.go
│   │   ├── proxy_controller.go
│   │   ├── lvyun_controller.go  # 绿云专用控制器
│   │   └── response.go
│   └── middleware/         # 中间件
│       └── apikey_auth.go  # API Key鉴权
├── service/                # 服务层
│   ├── config/             # 配置管理
│   │   └── env_config.go
│   ├── init.go             # 服务初始化
│   └── proxy/              # 代理服务
│       ├── interface.go    # 接口定义
│       ├── registry.go     # 注册中心
│       └── lvyun/          # 绿云客户端
│           ├── client.go   # 客户端实现
│           ├── models.go   # 数据模型
│           └── client_test.go # 单元测试
├── docs/                   # Swagger文档
└── scripts/
    └── start.sh            # 启动脚本
```

## 环境变量配置

### 基础配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `LISTEN_PORT` | 服务监听端口 | `80` |
| `BASE_CONTEXT` | 基础路径上下文(可选) | `` |
| `API_KEY` | API访问密钥 | `default-api-key-please-change-in-production` |

### 绿云接口配置

| 变量名 | 说明 | 是否必需 |
|--------|------|----------|
| `LVYUN_BASE_URL` | 绿云接口基础URL | 是 |
| `LVYUN_HOTEL_GROUP_CODE` | 集团代码 | 是 |
| `LVYUN_USER_CODE` | 用户名 | 是 |
| `LVYUN_PASSWORD` | 密码 | 是 |
| `LVYUN_APP_KEY` | AppKey | 是 |
| `LVYUN_APP_SECRET` | AppSecret | 是 |

## 快速开始

### 本地开发

1. **克隆项目**
```bash
cd /Users/liu/Work/go-project/dapr-platform/datahub-external
```

2. **安装依赖**
```bash
go mod tidy
```

3. **生成Swagger文档**
```bash
swag init --parseDependency --parseInternal --parseDepth 1
```

4. **配置环境变量**
```bash
export API_KEY=your-api-key
export LVYUN_BASE_URL=http://your-lvyun-server:port
export LVYUN_HOTEL_GROUP_CODE=MHZSQG
export LVYUN_USER_CODE=yjyMHZSQGKPI
export LVYUN_PASSWORD=yjyMHZSQGKPI
export LVYUN_APP_KEY=40130
export LVYUN_APP_SECRET=7009544c428b8118ac593d5d1f5118d9
```

5. **启动服务**
```bash
cd scripts
./start.sh
```

或者直接运行:
```bash
go run main.go
```

6. **访问Swagger文档**
```
http://localhost:8080/swagger/index.html
```

### Docker部署

1. **构建镜像**
```bash
docker build -t datahub-external:latest .
```

2. **运行容器**
```bash
docker run -d \
  -p 8080:80 \
  -e API_KEY=your-api-key \
  -e LVYUN_BASE_URL=http://your-lvyun-server:port \
  -e LVYUN_HOTEL_GROUP_CODE=MHZSQG \
  -e LVYUN_USER_CODE=yjyMHZSQGKPI \
  -e LVYUN_PASSWORD=yjyMHZSQGKPI \
  -e LVYUN_APP_KEY=40130 \
  -e LVYUN_APP_SECRET=7009544c428b8118ac593d5d1f5118d9 \
  --name datahub-external \
  datahub-external:latest
```

## API使用示例

### 认证

所有API请求(除健康检查外)都需要提供API Key，支持两种方式:

**方式1: X-API-Key Header**
```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/datasources
```

**方式2: Authorization Bearer**
```bash
curl -H "Authorization: Bearer your-api-key" http://localhost:8080/datasources
```

### 健康检查

```bash
# 基础健康检查
curl http://localhost:8080/health

# 就绪检查
curl http://localhost:8080/ready
```

### 列出所有数据源

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/datasources
```

响应示例:
```json
{
  "status": 200,
  "data": [
    {
      "name": "lvyun",
      "route_prefix": "/lvyun",
      "description": "绿云酒店管理系统接口",
      "status": "active"
    }
  ]
}
```

### 绿云接口调用示例

> **注意**: 所有日期参数格式为 `YYYY-MM-DD HH:MM:SS`，在URL中需要编码为 `YYYY-MM-DD%20HH:MM:SS`

**1. 获取绿云接口信息**
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8180/lvyun/info"
```

**2. 绿云健康检查**
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8180/lvyun/health"
```

**3. 预订单数据**
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8180/lvyun/reservations?hotel_code=HOTEL001&start_date=2025-10-20%2000:00:00&end_date=2025-10-21%2000:00:00"
```

响应示例:
```json
{
  "status": 200,
  "data": [
    {
      "Hotelcode": "HOTEL001",
      "HotelName": "示例酒店",
      "Id": "123456",
      "Name": "张三",
      "RmType": "标准间",
      "Rmno": "2001",
      "ArrDate": "2025-10-20",
      "DepDate": "2025-10-21",
      "RsvNo": "RSV001",
      "Mobile": "13800138000",
      "Sta": "R"
    }
  ]
}
```

**4. 登记单数据**
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8180/lvyun/registrations?hotel_code=HOTEL001&start_date=2025-10-20%2000:00:00&end_date=2025-10-21%2000:00:00"
```

**5. 结账单数据**
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8180/lvyun/checkouts?hotel_code=HOTEL001&biz_date=2025-10-20%2000:00:00"
```

**6. 营业报表数据**
```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8180/lvyun/business-report?hotel_code=HOTEL001&start_date=2022-10-20%2000:00:00&end_date=2022-10-21%2000:00:00"
```

## 从 datahub-service 调用

在 `datahub-service` 中可以通过HTTP接口调用:

```go
// 配置datahub-external地址
externalURL := "http://datahub-external:8180"
apiKey := "your-api-key"

// 调用绿云预订单接口
req, _ := http.NewRequest("GET", 
    externalURL+"/lvyun/reservations?hotel_code=HOTEL001&start_date=2025-10-20%2000:00:00&end_date=2025-10-21%2000:00:00", nil)
req.Header.Set("X-API-Key", apiKey)

resp, err := http.DefaultClient.Do(req)
// 处理响应...
```

## 扩展新的第三方接口

添加新的第三方接口客户端遵循以下步骤:

### 1. 创建数据模型

在 `service/proxy/yourclient/models.go` 中定义请求和响应结构体:

```go
package yourclient

// YourRequest 请求参数
type YourRequest struct {
    Param1 string `json:"param1"`
    Param2 int    `json:"param2"`
}

// YourResponse 响应结果
type YourResponse struct {
    Data1 string `json:"data1"`
    Data2 int    `json:"data2"`
}
```

### 2. 实现客户端

在 `service/proxy/yourclient/client.go` 中实现 `ThirdPartyClient` 接口:

```go
type ThirdPartyClient interface {
    GetName() string
    GetRoutePrefix() string
    GetDescription() string
    HandleRequest(ctx context.Context, path string, params map[string]string) (interface{}, error)
    HealthCheck(ctx context.Context) error
    GetStatus() string
}
```

### 3. 创建专用控制器

在 `api/controllers/yourclient_controller.go` 中创建控制器，定义具体的API接口，避免使用 `map[string]interface{}`：

```go
// GetData 获取数据
// @Summary 获取数据
// @Description 详细描述
// @Tags 您的客户端
// @Param param1 query string true "参数1"
// @Success 200 {object} APIResponse{data=[]yourclient.YourResponse}
// @Security ApiKeyAuth
// @Router /yourclient/data [get]
func (c *YourClientController) GetData(w http.ResponseWriter, r *http.Request) {
    // 具体实现
}
```

### 4. 注册客户端

在 `service/init.go` 中注册:

```go
newClient := yourclient.NewYourClient(config)
proxy.GetGlobalRegistry().Register(newClient)
```

### 5. 注册路由

在 `api/routes.go` 中添加路由:

```go
yourClientController := controllers.NewYourClientController(proxy.GetGlobalRegistry())
r.Route("/yourclient", func(r chi.Router) {
    r.Get("/health", yourClientController.HealthCheck)
    r.Get("/data", yourClientController.GetData)
})
```

### 6. 生成Swagger文档

```bash
swag init --parseDependency --parseInternal --parseDepth 1
```

## 技术栈

- **Go**: 1.23.1
- **Web框架**: Chi v5
- **文档**: Swagger/OpenAPI
- **容器化**: Docker
- **微服务框架**: Dapr

## 设计原则

- ✅ **KISS原则**: 保持简单，无数据库依赖
- ✅ **YAGNI原则**: 只实现必需功能
- ✅ **DRY原则**: 复用公共库和模式
- ✅ **单一职责**: 每个组件职责清晰
- ✅ **依赖倒置**: 面向接口编程

## 注意事项

1. **生产环境务必修改API_KEY**，不要使用默认值
2. 绿云接口的Session会自动刷新(每11小时)
3. 所有请求都会经过API Key认证(除白名单路径)
4. 建议配置合适的超时时间和重试策略

## 许可证

与主项目保持一致

## 联系方式

如有问题，请联系项目维护团队


# datahub-external
