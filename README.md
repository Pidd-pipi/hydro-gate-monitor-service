# Hydro Gate Monitor Service

Hydroelectric gate observation service with an in-memory operational store. It exposes gate observations and allows an operator to acknowledge a non-normal alert.

## Layout

`config`, `domain`, `store`, `validation`, `health`, and `api` are separate service layers. The embedded `web` directory is served by the binary.

## Run and test

The service reads `PORT` and defaults to `8080`.

```text
cd backend && go build ./...
cd backend && go test ./...
cd backend && PORT=8080 go run .
```

`GET /healthz` returns service health. `GET /api/gates` returns gate observations. `POST /api/gates/gate-01/ack` accepts `{"note":"Operator reviewed spillway reading"}` and acknowledges the alert. Invalid or empty notes return `400`.

## Verification

Verified from `backend/` with `gofmt -w .`, `go build ./...`, and `go test ./...`. A live service on `PORT=18180` returned `200` for `/healthz`, `/api/gates`, `/`, and `/app.js`; the valid acknowledgement returned `200` and an empty note returned `400`. The service was stopped in the same smoke operation after checks.

## Engineering Notes

水闸运行流程代码按领域模型、校验、状态转换、并发安全存储、审计事件和 HTTP 生命周期分层。请求会保留请求标识并经过恢复与超时保护；状态写入使用版本校验，错误通过可识别的领域错误返回。

除现有接口回归测试外，项目还保留可复用的分页、过滤、策略、工作流和运行健康能力，便于后续扩展而不把业务规则堆积到处理器中。

## Enterprise Layout

```text
.
├── backend/       # Go module, source, tests, and embedded web assets
├── database/      # persistence documentation and future schemas
├── output/        # verification records
├── prompt.txt
└── runtime_smoke.json
```

Run `cd backend && go test ./...`, `cd backend && go build ./...`, or `cd backend && PORT=8080 go run .`. The health check is `GET /healthz`; the main API endpoints are `GET /api/gates` and `POST /api/gates/{id}/ack`.
