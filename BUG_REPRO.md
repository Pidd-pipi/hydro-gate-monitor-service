# Bug Reproduction — hydro-gate-monitor-service-010

## Bug 是什么
静态资源缓存与健康检查的并发问题：包级 assetCache map 无锁并发读写，
页面并发打开直接崩（concurrent map read and map write）；健康检查计数器
非原子自增并发竞态；健康检查不校验请求方法，POST 也返回 ok。

## 如何触发
- 多个并发请求同时打开页面（冷缓存首次加载）。
- 并发打健康检查；用 POST 请求 /healthz。

## 真实错误信息
```
fatal error: concurrent map read and map write

goroutine 21 [running]:
example.com/hydro-gate-monitor-service/api.NewRouter.func1()
	.../backend/api/router.go:26 +0xc0
net/http.HandlerFunc.ServeHTTP()
	.../src/net/http/server.go:2322 +0x38
```
