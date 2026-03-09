# CloudOps

CloudOps 是一个基于 Go 的云原生运维管理平台，提供 Kubernetes 多集群管理、Prometheus 告警监控、RBAC 权限控制、工单系统和 AI 智能运维（RAG）能力。

## 技术栈

| 类别 | 技术 |
|------|------|
| 后端语言 | Go 1.24+ |
| HTTP 框架 | Gin |
| ORM | GORM + MySQL |
| 缓存 | Redis |
| 权限 | Casbin（RBAC）+ JWT |
| K8s 客户端 | client-go |
| 依赖注入 | Wire（google/wire） |
| 日志 | Uber zap |
| API 文档 | Swagger（swaggo/swag） |
| AI 模块 | OpenAI 兼容 API + 内存向量检索（RAG） |

## 功能模块

### K8s 多集群管理
- 集群注册、连接状态管理（kubeconfig 动态加载）
- 多集群 clientset 池（`sync.Map` 维护）
- Namespace / Deployment / Pod / Service / ConfigMap 增删查
- Deployment 扩缩容、Pod 日志实时获取

### Prometheus 告警
- Scrape Pool 管理（绑定 Prometheus 和 AlertManager 地址）
- 告警规则（PromQL + 持续时间 + 严重级别）CRUD
- AlertManager Webhook 接收告警事件，异步通知（钉钉 / 邮件 / Webhook）
- 发送组（Send Group）灵活配置通知渠道

### RBAC 权限系统
- User / Role / Menu 三层模型
- Casbin enforcer（文件模型 + DB 策略）
- JWT 双 Token（短期 access + 长期 refresh）
- 管理员（AccountType=2）直接放行

### 工单系统
- 工单模板（自定义 JSON Schema 表单）
- 工单实例生命周期状态机：
  ```
  pending → in_progress → resolved → closed
                        ↘ rejected ← reopen ←
  ```
- 流转记录审计、评论协作

### AI 智能运维（RAG）
- 知识库文档管理（MySQL 持久化 + embedding JSON 字段）
- 启动时将 embedding 加载到内存向量存储
- 余弦相似度 Top-K 检索，构建 RAG prompt
- 对话接口兼容 OpenAI API（可替换任意兼容服务）

## 目录结构

```
CloudOps/
├── main.go                     # 入口：加载配置 → Wire 初始化 → 启动
├── generate.go                 # //go:generate wire ./pkg/di/
├── config/
│   ├── config.yaml             # 服务配置
│   └── rbac_model.conf         # Casbin 模型
├── docs/                       # Swagger 文档（swag init 生成）
├── internal/
│   ├── model/                  # 数据模型 + 请求结构体
│   ├── middleware/             # JWT / Casbin / CORS / Logger
│   ├── k8s/
│   │   ├── api/                # Cluster / Namespace / Deployment / Pod / Service / ConfigMap
│   │   ├── service/
│   │   ├── dao/
│   │   └── client/             # K8sClient：sync.Map 多集群管理
│   ├── prometheus/
│   │   ├── api/                # ScrapePool / AlertRule / AlertEvent / SendGroup
│   │   ├── service/
│   │   ├── dao/
│   │   └── notify/             # DingTalk / Email / Webhook 通知实现
│   ├── system/
│   │   ├── api/                # User / Role / Menu
│   │   ├── service/
│   │   └── dao/
│   ├── workorder/
│   │   ├── api/                # Template / Instance / Flow / Comment
│   │   ├── service/            # flow_service.go 实现状态机
│   │   └── dao/
│   ├── ai/
│   │   ├── api/                # /chat  /knowledge
│   │   ├── service/            # llm_service.go + rag_service.go
│   │   ├── dao/
│   │   └── vector/             # 内存余弦相似度向量检索
│   └── startup/                # 启动任务：加载向量知识库
└── pkg/
    ├── base/                   # 统一响应格式 + HandleRequest
    ├── jwt/                    # JWT Handler（双 Token）
    └── di/                     # Wire 依赖注入（wire.go + wire_gen.go）
```

## 快速开始

### 前置依赖

- Go 1.24+
- MySQL 8.0+
- Redis 6+

### 配置

```bash
cp .env.example .env
# 编辑 config/config.yaml，修改数据库、Redis 连接地址
```

`config/config.yaml` 关键配置项：

```yaml
server:
  port: "8080"

mysql:
  addr: "root:root@tcp(localhost:3306)/cloudops?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "localhost:6379"

ai:
  openai_api_key: "sk-xxx"
  openai_base_url: "https://api.openai.com/v1"
  chat_model: "gpt-4o-mini"
  embedding_model: "text-embedding-3-small"
```

### 运行

```bash
# 安装依赖
go mod tidy

# 启动（首次运行自动建表）
go run main.go
```

### API 文档

```bash
# 安装 swag（如未安装）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g main.go -o docs/

# 启动后访问
open http://localhost:8080/swagger/index.html
```

### 重新生成 Wire DI

```bash
go install github.com/google/wire/cmd/wire@latest
go generate ./...
```

## API 路由概览

| 模块 | 路由前缀 |
|------|---------|
| 认证 | `POST /api/auth/login` `POST /api/auth/logout` |
| 用户 | `GET/POST/PUT/DELETE /api/users` |
| 角色 | `GET/POST/PUT/DELETE /api/roles` |
| 菜单 | `GET/POST/PUT/DELETE /api/menus` |
| K8s 集群 | `GET/POST/PUT/DELETE /api/k8s/clusters` |
| Namespace | `GET/POST/DELETE /api/k8s/namespaces` |
| Deployment | `GET/POST/DELETE /api/k8s/deployments` `POST /api/k8s/deployments/scale` |
| Pod | `GET/DELETE /api/k8s/pods` `GET /api/k8s/pods/logs` |
| Scrape Pool | `GET/POST/PUT/DELETE /api/prometheus/scrape-pools` |
| 告警规则 | `GET/POST/PUT/DELETE /api/prometheus/alert-rules` |
| 告警事件 | `GET /api/prometheus/alert-events` `POST /api/prometheus/alert-events/receive` |
| 发送组 | `GET/POST/PUT/DELETE /api/prometheus/send-groups` |
| 工单模板 | `GET/POST/PUT/DELETE /api/workorder/templates` |
| 工单实例 | `GET/POST /api/workorder/instances` |
| 工单流转 | `GET /api/workorder/instances/:id/flows` `POST .../action` |
| AI 对话 | `POST /api/ai/chat` |
| 知识库 | `GET/POST/DELETE /api/ai/knowledge` |

## 数据库表

| 表名 | 说明 |
|------|------|
| users / roles / menus | 系统 RBAC |
| user_roles / role_menus | 多对多关联 |
| k8s_clusters | K8s 集群注册信息 |
| prom_scrape_pools | Prometheus 实例池 |
| prom_alert_rules | 告警规则 |
| prom_alert_events | 告警事件记录 |
| prom_send_groups | 通知发送组 |
| wo_templates | 工单模板 |
| wo_instances | 工单实例 |
| wo_flows | 流转记录 |
| wo_comments | 评论 |
| ai_knowledge_chunks | AI 知识库（含 embedding） |

## License

MIT
