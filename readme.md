
# honey 是什么

honey是一个开源的基于 [zapp](https://github.com/zly-app/zapp) 的日志收集处理工具.
目的是收集任何基于 `zapp` 开发的项目的日志并转储到任何地方.

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

# 结构图

待补充...

# 架构图

待补充...

# 如何运行

# 配置

待补充...

# zapp日志收集插件

待补充...
