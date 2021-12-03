## 模块8作业1

- 优雅启动： 通过depoyment 中的 initialDelaySeconds 来延迟检测时间
- 优雅终止: 在代码中捕获SIGTERM信号处理优雅停止，没有用到kubernetes的 prestop hook
- Qos 保证：使用 Burstable pods,  requests < limits
- 探活: 是kubernetes去做这个事情的，体验在 rediness probe 和 liveness probe
- 日志等级: 通过配置文件控制
- 配置与代码分离: 通过 configmap 以及代码中检测配置文件变更自动变更配置

代码链接: https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/httpserver/main.go
