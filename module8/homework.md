## 模块8作业

- 优雅启动： 通过depoyment 中的 initialDelaySeconds 来延迟检测时间
- 优雅终止: 在代码中捕获SIGTERM信号处理优雅停止，没有用到kubernetes的 prestop hook
- Qos 保证：使用 Burstable pods,  requests < limits
- 探活: 是kubernetes去做这个事情的，体验在 rediness probe 和 liveness probe
- 日志等级: 通过配置文件控制
- 配置与代码分离: 通过 configmap 以及代码中检测配置文件变更自动变更配置
- 应用高可用:
   1. 尽量保证应用无状态
   2. 使用 Depoyment 部署应用 ,并且 replicas 设置成多个实例
   3. 应用本身需要做到优雅终止，已经进来的流量要处理完成再终止
   4. 配置好探活
   5. 配置好资源，requests 和 limit


代码链接: https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/httpserver/main.go
kubernetes配置文件链接:  https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/httpserver/depoyment

