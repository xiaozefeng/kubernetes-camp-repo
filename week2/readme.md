
## 需求
1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回 200

## 作业链接：
#### net/http 实现
https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/week2/homework/nethttp/main.go

### net 包实现
https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/week2/homework/tcp/main.go