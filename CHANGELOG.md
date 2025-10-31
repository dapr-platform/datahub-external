# 更新日志

## [优化版本] - 2025-10-30

### 新增
- ✨ 添加 `BASE_CONTEXT` 支持，允许在指定路径下挂载所有路由
- ✨ 创建绿云专用控制器 `api/controllers/lvyun_controller.go`
- ✨ 创建绿云数据模型 `service/proxy/lvyun/models.go`
- ✨ 创建 API 测试脚本 `scripts/test_api.sh`
- ✨ 创建优化总结文档 `OPTIMIZATION_SUMMARY.md`
- ✨ 创建更新日志 `CHANGELOG.md`

### 改进
- 🔧 完善 Swagger 文档配置，添加 API Key 安全定义
- 🔧 优化接口路由结构，使用专用路由代替通用代理
- 🔧 更新单元测试使用真实 URL (http://182.151.37.189:8102)
- 🔧 改进错误处理和参数验证
- 📝 更新 README.md，完善文档说明
- 📝 更新 env.example，添加 BASE_CONTEXT 配置

### 变更
- 💥 **重要**: API 路由结构变更
  - 旧: `/datasource/lvyun/{path}`
  - 新: `/lvyun/{specific-endpoint}`
- 🔄 绿云路由前缀从 `/datasource/lvyun` 改为 `/lvyun`
- 🔄 简化 `proxy_controller.go`，只保留 `ListDataSources` 方法

### 移除
- ❌ 删除通用代理接口 `ProxyRequest`
- ❌ 删除 `proxy_controller.go` 中的 `HealthCheck` 方法（移到专用控制器）

### 修复
- 🐛 修复导入未使用的问题
- 🐛 统一使用 snake_case 参数名（如 `hotel_code`）

### 新的 API 端点

#### 数据源管理
```
GET /datasources              # 列出所有数据源
```

#### 绿云酒店接口
```
GET /lvyun/info               # 获取接口信息
GET /lvyun/health             # 健康检查
GET /lvyun/room-total         # 实时房情数据 (需要 hotel_code 参数)
GET /lvyun/room-rank          # 酒店排名
GET /lvyun/room-adr           # 平均房价走势 (需要 hotel_code 参数)
GET /lvyun/room-detail        # 营收情况详情 (需要 hotel_code 参数)
GET /lvyun/room-detail-list   # 营收情况列表 (需要 hotel_code 参数)
```

### 环境变量
```bash
# 新增
BASE_CONTEXT=/api/v1  # 可选，用于指定基础路径

# 已有
LISTEN_PORT=8180
API_KEY=your-api-key
LVYUN_BASE_URL=http://182.151.37.189:8102
LVYUN_HOTEL_GROUP_CODE=MHZSQG
LVYUN_USER_CODE=yjyMHZSQGKPI
LVYUN_PASSWORD=yjyMHZSQGKPI
LVYUN_APP_KEY=40130
LVYUN_APP_SECRET=7009544c428b8118ac593d5d1f5118d9
```

### 迁移指南

如果你正在使用旧版本的 API，需要更新调用方式：

#### 旧版本调用方式
```bash
# 实时房情
curl -H "X-API-Key: key" \
  "http://localhost:8180/datasource/lvyun/room-total?hotelCode=HOTEL001"

# 健康检查
curl -H "X-API-Key: key" \
  "http://localhost:8180/datasource/lvyun/health"
```

#### 新版本调用方式
```bash
# 实时房情
curl -H "X-API-Key: key" \
  "http://localhost:8180/lvyun/room-total?hotel_code=HOTEL001"

# 健康检查
curl -H "X-API-Key: key" \
  "http://localhost:8180/lvyun/health"
```

**关键变更**:
1. 路径从 `/datasource/lvyun/` 改为 `/lvyun/`
2. 参数名从驼峰式 `hotelCode` 改为下划线式 `hotel_code`
3. 每个接口都有明确的端点，不再使用通配符路由

### 技术债务清理
- ✅ 去除了过度通用的设计
- ✅ 提高了类型安全性
- ✅ 改进了代码可读性
- ✅ 增强了 API 文档完整性

### 测试覆盖
- ✅ 单元测试：签名计算、接口实现
- ✅ 集成测试：真实环境测试
- ✅ API 测试脚本：自动化测试所有接口

### 性能
- 无明显性能变化
- 代码简化可能略微提升性能

### 兼容性
- ⚠️ **破坏性变更**: API 路由结构变更
- ⚠️ **破坏性变更**: 参数命名规范变更
- ✅ 环境变量向后兼容
- ✅ 新增 BASE_CONTEXT 为可选配置

### 文档
- 📚 完善了 README.md
- 📚 添加了 OPTIMIZATION_SUMMARY.md
- 📚 更新了 API 使用示例
- 📚 完善了扩展指南

### 下一步计划
- [ ] 添加更多第三方接口客户端
- [ ] 实现请求限流
- [ ] 添加请求日志记录
- [ ] 实现缓存机制
- [ ] 添加监控指标

---

## 贡献者
- 优化和重构：AI Assistant
- 项目维护：开发团队

## 反馈
如有问题或建议，请联系项目维护团队。

