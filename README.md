# 兰州水环境治理平台

这是一个面向河湖断面监测、上下游联防联控、横向生态补偿、固体废物转移和生态项目验收的 Go 后端。系统用 SQLite 保存真实业务状态，支持迁移、重启恢复、事务回滚、乐观并发、后台重试和审计。

## 运行

```bash
GOTOOLCHAIN=local go test ./... -count=1
GOTOOLCHAIN=local go test -race ./... -count=1
GOTOOLCHAIN=local go vet ./...
GOTOOLCHAIN=local go build ./...
go run ./cmd/server
```

健康检查为 `/healthz`，就绪检查为 `/readyz`。默认数据库位于 `./data/lanzhou.db`。
