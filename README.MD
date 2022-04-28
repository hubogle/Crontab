# Crontab 分布式调度项目

`master` 负责任务创建，页面接口对接
`worker` 负责任务的调度，以及任务的执行
`jobs` 一些需要被调度的任务执行逻辑

# 项目技术

## 常用包

1. `gin`
2. `viper`
3. `gorm`

## 三方工具

1. `Consul` 服务注册发现
2. `Loki` 日志监控
3. `Grafana` 数据可视化
4. `Prometheus` 数据监控
5. `Nacos` 配置文件管理
