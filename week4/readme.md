## 构建本地镜像。
编写好源代码 + Dockerfile 在当前目录执行
docker build -t image_name:tag .

## 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化（请思考有哪些最佳实践可以引入到 Dockerfile 中来）。
https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/httpserver/main.go
https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/httpserver/Makefile
https://github.com/xiaozefeng/kubernetes-camp-repo/blob/master/httpserver/Dockerfile


## 将镜像推送至 Docker 官方镜像仓库。
1. docker login 登录
2. docker push image_name:tag


## 通过 Docker 命令本地启动 httpserver。
docker run -d  --name httpserver -p 9000:9000 httpserver:v1.0

## 通过 nsenter 进入容器查看 IP 配置
1. 获取容器在宿主机上的PID
```bash
docker inspect -f {{.State.Pid}} 容器id
```
2. 是用nsenter +pid 获取进程 网络配置
```bash
nsenter -t $PID -n ip a
```