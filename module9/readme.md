# 模块9

## QOS Class
BestEffort  没有资源 request 和 limit
Burstable  资源限制的 request < limit
Guaranteed 资源限制的 request = limit

驱逐优先级从高到低
BestEffort > Burstable > Guaranteed


## 节点错误上报
helm 包管理
node-problem-detector  是一个daemonset， 只负责上报错误，在node的event中加入一个warning的message
需要自己写控制器去定义行为



## debug 
1. systemd 拉起的进程日志 journalctl -auf kubelet
2. pod 日志  k logs -f <podname>
3. pod 某个容器的日志 k logs -c <containerName> <podName>

## 监控
普罗米修斯
loki 日志系统
grafna 


