### 创建 eureka-server的工程
配置文件看 application.yaml

### 构建 Docker 镜像
看 Dockerfile

### 创建service 对象
由于需要使用headless service 和外部访问 eureka-server 需要两个service ， 看 service.yaml


### 由于要保持稳定的网络标识，需要使用statefulset
看 statefulset.yaml