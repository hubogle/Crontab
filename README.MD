# Crontab 分布式调度项目

# 目录结构

``shell
├── app
│ ├── Jobs
│ ├── master
│ └── worker    
├── cmd
│ └── generdate # 通过 gorm-gen 生成 Model
├── go.mod
├── README.MD
├── test
│ └── api # http api 接口测试
└── util
├── db # Db 工具包
└── validation # validator 共用方法
``

# 项目介绍

`master` 负责任务创建，页面接口对接
`worker` 负责任务的调度，以及任务的执行
`jobs` 一些需要被调度的任务执行逻辑

## master

1. `Gin` 提供 `Job` 的增删改查功能。
2. 通过 `consul` 操作 `cron/jobs/` 前缀 `key`，通过 `worker` 节点进行 `watch` 任务的变动。

## worker

1. 将服务注册到 `consul` 负责服务发现功能；
2. 使用 `consul` 的 `KV` 将每个任务的配置信息进行保存；
3. 利用 `consul` 的 `watch` 机制对某个目录下所有的 `key` 进行监控；
4. 监听前缀目录`cron/jobs/`所有任务的状态，将任务放入到 `scheduler` 中进行调度；
5. 使用 `consul` 中 `Session` 配合 `Lock` 实现分布式锁，防止任务并发执行；
   1. 首先申请一个 `Session`，设置 `TTL` 时间为 `20s`。
   2. 基于该 `Session` 对 `KV` 调用 `Acquire` 进行锁的获取。
   3. 若其中一个节点获取到了锁，则进行该任务的执行。
   4. 启动一个 `goroutine` 对该 `Session` 进行 `Renew` 操作，每隔 `20s` 进行一次 `Renew` 操作。
   5. 当任务执行完成后，对该 `Session` 的 `KV` 调用 `Release` 方法释放锁（也可等 `TTL` 到期自动删除）。
   6. 停止掉 `goroutine` 对 `Session` 的 `Renew` 操作。

确保所有节点都有机会执行任务，引入随机睡眠 `0~1s`。

# 基于GIN的一些封装

1. `pkg.core.context` 使用 `sync.Pool` 对 `*gin.Context` 做了封装，使得 `*gin.Context` 可以被复用。

## 常用包

1. `gin`
2. `viper`
3. `gorm`
5. [Gorm-Gen](https://github.com/go-gorm/gen/blob/master/README.ZH_CN.md)

## 三方工具

1. `Consul` 服务注册发现
2. `Loki` 日志监控
3. `Grafana` 数据可视化
4. `Prometheus` 数据监控
5. `Nacos` 配置文件管理
