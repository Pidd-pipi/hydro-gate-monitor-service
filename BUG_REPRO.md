# Bug Reproduction — hydro-gate-monitor-service-003

## Bug 是什么
错误分类与响应状态码错乱：错误包装丢失哨兵链，分类比较用 `==` 无法识别
包装后的错误，全部塌陷成 internal；确认不存在的闸门返回 409 而不是 404。

## 如何触发
- 触发"记录不存在 / revision 冲突 / 策略拒绝"等错误并经过包装传递。
- POST /api/gates/{不存在的id}/ack。

## 真实错误信息
`errors.Is(err, ErrOpsNotFound)` 对包装后的错误返回 false；
不存在的闸门确认接口返回 `409 conflict`（应为 404）。
