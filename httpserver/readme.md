## module 10 作业

1. 为 HTTPServer 添加 0-2 秒的随机延时
    ``` go
   delay := randInt(10, 2000)
    time.Sleep(time.Millisecond * time.Duration(delay))
   ```
2. 为 HTTPServer 项目添加延时 Metric
    ``` go
   mux.Handle("/metrics", promhttp.Handler())
   ```
   
3. 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置
   1. 集群安装 loki + prometheus + grafana 一整套
   2. 将http-server 打包成 helm的package
   2. 使用helm 安装 http-server 的package
   2. 使用 prometheus 的 dashboard
4. 从 Prometheus 界面中查询延时指标数据
（可选）创建一个 Grafana Dashboard 展现延时分配情况
