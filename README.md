# 多渠道告警编排与通知路由引擎

这是一个纯Go告警接入服务：接收监控、IoT和应用事件，按规则创建或聚合告警，支持静默窗口、确认、恢复，并把通知任务交给可替换的Webhook、邮件和短信适配器。默认使用线程安全内存仓储，接口边界和SQL迁移已为PostgreSQL、Redis与消息队列替换预留。

## 快速启动

```sh
cd project-02-alert-orchestration-engine
go test ./...
go run ./cmd/alertd
```

默认监听 `0.0.0.0:8092`。可通过 `ALERT_PORT`、`ALERT_HOST`、`ALERT_ENV`、`ALERT_SHUTDOWN_SECONDS`、`ALERT_LOG_LEVEL` 覆盖配置，也可设置 `ALERT_CONFIG=configs/config.yaml` 读取YAML。

## 核心流程

```sh
curl -sS -X POST http://127.0.0.1:8092/api/v1/rules -H 'Content-Type: application/json' -d '{"name":"CPU critical","eventName":"cpu.high","enabled":true,"severity":"critical","channels":["webhook","email"]}'
curl -sS -X POST http://127.0.0.1:8092/api/v1/events -H 'Content-Type: application/json' -d '{"source":"node-1","name":"cpu.high","severity":"critical","description":"CPU above limit","labels":{"region":"jp"},"value":98}'
curl -sS 'http://127.0.0.1:8092/api/v1/alerts?status=open'
curl -sS 'http://127.0.0.1:8092/api/v1/notifications'
curl -sS -X POST http://127.0.0.1:8092/api/v1/alerts/ALT_ID/acknowledge
curl -sS -X POST http://127.0.0.1:8092/api/v1/alerts/ALT_ID/resolve
```

## 目录结构

`cmd/alertd`为服务入口；`internal/domain`按event、rule、alert、silence、schedule、notification划分；`internal/application`包含用例和端口接口；`internal/adapter/http`包含REST接口及中间件；`internal/infrastructure`包含内存仓储、通知、队列、日志和时钟适配器。`api/openapi.yaml`为接口契约，`migrations`包含up/down迁移，`deploy`包含Docker文件。

## 运行时端点

`GET /healthz`、`GET /readyz`、`GET /metrics`，以及事件、规则、告警、静默窗口、通知记录API。服务使用请求ID、访问日志、统一JSON错误、超时、CORS和优雅关闭。

## 质量约束

非测试Go源码目标为至少2000行，统计排除测试、生成代码、依赖目录、锁文件和构建产物。领域层不依赖HTTP或数据库，外部依赖通过接口和构造函数注入，所有处理链路传递context并包装底层错误。
