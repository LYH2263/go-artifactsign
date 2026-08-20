# go-artifactsign

制品摘要、HMAC 签名/验签、信任根与吊销清单库；配套 `signd` 管理服务与静态页。

## Build / Test

```bash
go build ./...
go test ./... -count=1
```

## Daemon

```bash
go run ./cmd/signd -addr :8099 -web web
```

Management page: http://127.0.0.1:8099/
