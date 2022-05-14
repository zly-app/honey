
<!-- TOC -->

- [honey 是什么](#honey-是什么)
- [关键概念](#关键概念)
  - [zapp](#zapp)
  - [input](#input)
  - [output](#output)
  - [compress](#compress)
  - [serializer](#serializer)
- [结构图](#结构图)
- [架构图](#架构图)
- [如何运行](#如何运行)
  - [本地编译运行](#本地编译运行)
  - [在 docker 运行](#在-docker-运行)
  - [使用 docker-compose](#使用-docker-compose)
- [配置](#配置)
- [zapp日志收集插件](#zapp日志收集插件)

<!-- /TOC -->

---

# honey 是什么

honey是一个开源的基于 [zapp](https://github.com/zly-app/zapp) 的日志收集处理工具. 目的是收集任何基于 `zapp` 开发的项目的日志并转储到任何地方.

# 关键概念

在深入了解 `honey` 之前, 最好熟悉该服务的一些关键概念, 我们再下面提供了我们定义的部分术语描述.

## zapp

`zapp` 是一个用于快速构建项目的基础库, 点 [这里](https://github.com/zly-app/zapp) 了解更多.

## input

`input` 从其它任何地方接收日志, 比如 `honey-http`, `kafka`, `pulasr` 等

## output

`output` 将日志处理并存放到用于日志分析的地方, 如 `clickhouse`, `elasticsearch`, `loki`, `Influx` 等

## compress

压缩器, 目前支持 `raw`, `gzip`, `zstd`

## serializer

序列化器, 目前支持 `json`, `msgpack`

# 结构图

待补充...

# 架构图

待补充...

# 如何运行

## 本地编译运行

```sh
git clone https://github.com/zly-app/honey.git
cd honey
go run .
```

## 在 docker 运行

```sh
docker run zlyuan/honey:v0.2.0
```

## 使用 docker-compose

```sh
wget https://github.com/zly-app/honey/blob/master/docker-compose.yml
docker-compose up -d
```

# 配置

默认配置文件路径 `./configs/default.toml`, 默认配置文件可以不存在. 使用 `-c` 以指定配置文件启动, 此时配置文件必须存在.

详细配置文件书写参考[这里](./configs/default.toml)

示例:

```toml
# honey配置项
[services.honey]
# 输入设备列表, 多个输入设备用半角逗号`,`分隔, 目前支持的输入设备: http
Inputs = 'http'
# 输出设备列表, 多个输出设备用半角逗号`,`分隔, 目前支持的输出设备: std, honey-http, loki-http
Outputs = 'std'

# http输入器插件配置项
[input.http]
# 监听地址, 示例: :8080
#Bind = ':8080'
# 推送
#PushPath = '/path'
```

# zapp日志收集插件

转到 [这里](https://github.com/zly-app/plugin/tree/master/honey)
