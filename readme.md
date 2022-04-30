# honey 是什么

honey是一个开源的基于 [zapp](https://github.com/zly-app/zapp) 的日志收集处理工具. 目的是收集任何基于 `zapp` 开发的项目的日志并转储到任何地方.

# 关键概念

在深入了解 `honey` 之前, 最好熟悉该服务的一些关键概念, 我们再下面提供了我们定义的部分术语描述.

## zapp

`zapp` 是一个用于快速构建项目的基础库, 点 [这里](https://github.com/zly-app/zapp) 了解更多.

## reporter

`reporter` 将收集到的日志数据上报到 `honey` 或别的任何地方, 比如 `kafka`, `pulsar` 等

## input

`input` 接收`reporter`发送的日志数据, 也可以从其它任何地方接收日志, 比如 `kafka`, `pulasr` 等

## output

`output` 将日志处理并存放到用于日志分析的地方, 如 `clickhouse`, `elasticsearch`, `loki`, `Influx` 等

## compress

压缩器, 目前支持 `raw`, `zstd`

## serializer

序列化器, 目前支持 `json`, `msgpack`

# 结构图

待补充...

# 架构图

待补充...

# 如何运行

`go run .`

# 配置

默认配置文件路径 `./configs/default.toml`, 默认配置文件可以不存在. 使用 `-c` 以指定配置文件启动, 此时配置文件必须存在.

配置文件书写参考[这里](./configs/default.toml)

---

# zapp日志收集插件

用于收集zapp项目的日志

## 示例

```go
package main

import (
	"github.com/zly-app/honey/zapp_plugin/honey"
	"github.com/zly-app/zapp"
)

func main() {
	app := zapp.NewApp("test",
		honey.WithPlugin(), // 启用日志收集插件
	)
	defer app.Exit()

	app.Run()
}
```

## 配置

```toml
# honey日志收集插件配置项
[plugins.honey]
# 上报时标示的环境名
Env = 'dev'
# 上报时标示的服务名, 如果为空则使用app名
Service = ''
# 上报时标示的实例名, 如果为空则使用本地ip
Instance = ''
# 停止原有的日志输出
StopLogOutput = true
# 日志批次大小, 累计达到这个大小立即上报一次日志, 不用等待时间
LogBatchSize = 10000
# 自动旋转时间(秒), 如果没有达到累计上报批次大小, 在指定时间后也会立即上报
AutoRotateTime = 3
# 最大旋转线程数, 表示同时允许多少批次发送到输出设备
MaxRotateThreadNum = 10
# 上报者, 支持 std, http
Reports = 'std'
```
