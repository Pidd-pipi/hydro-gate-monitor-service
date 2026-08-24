# Bug Reproduction — hydro-gate-monitor-service-005

## Bug 是什么
分页查询与状态机问题：零值分页参数算出负的切片起点，搜索直接越界 panic；
状态机非法迁移失败仍把脏记录写进历史；读取的历史快照与内部共享底层数组，
外部修改会污染库内历史。

## 如何触发
- 不带参数（Page 为 0）调用分页搜索。
- 触发一次非法状态迁移后读取历史；修改读到的历史快照再读库内历史。

## 真实错误信息
```
panic: runtime error: slice bounds out of range [-25:] [recovered, repanicked]
goroutine 5 [running]:
example.com/hydro-gate-monitor-service.(*OpsService).Search()
	.../backend/ops_service.go:68 +0x330
```
