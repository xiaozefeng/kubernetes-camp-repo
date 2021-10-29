## Docker
### Installation
[centos](https://docs.docker.com/engine/install/centos/)

### Core Concept

#### Image
> 镜像：打包了 软件以及软件所需要的操作系统 (但是不包含 Linux Kernel)

##### 镜像有什么好处？ 
1. 因为image本质是一个压缩包，所以可以很方便通过网络传播
2. 因为image打包的是软件已经软件依赖的操作系统，所以给予了一个宝贵的能力: 保证了本地环境与云端环境的高度一致 (除了内核，内核可能不一样)


##### 如何自定义镜像？
Docker 定义了一套打包的机制
1. 通过启动一个容器，进入容器，修改内容，保存之后，再打包成镜像
1. 通过 Dockerfile （推荐）


#### Repository
镜像仓库：存放镜像的地方，比如一个ubuntu镜像仓库，可以存放 ubuntu的多个版本，如 14.06， 16.06


#### Registry
注册服务器：存放镜像仓库地方

##### Image ，Repository， Registry的关系梳理
Registry 包含 Repository  包含  Image 
在一个注册服务器上 有很多 镜像仓库， 每个镜像仓库有多个版本的 Docker镜像

#### Container
容器：一个被namespace隔离起来，Cgroup 管控起来的，拥有独立文件系统的 进程
容器是镜像的动态视图

#### Container Network
容器网络

#### Container Storage
容器存储

### Docker 常用命令
docker images
docker ps
docker network ls 
docker inspect 
docker volume ls 
docker rm 
docker rmi
docker build 
docker run 
docker exec 


## 底层技术
- Namespace 资源隔离 (进程，网络，文件系统，用户，UTS， IPC通信)
- CGroups 资源限制 (CPU， Memory， IO)
- UnionFS  文件系统


### Namespace 隔离
包括 进程隔离，网络隔离，文件系统隔离, 用户隔离, IPC（interprocess communication）隔离
- pid namespace
- mnt namespace
- net namespace
- ipc namespace
- uts namespace
- user namespace

### 常用命令
nsenter 
lsns 

/proc 下面是反应系统的状态


### CGroups 资源控制
资源管控


### UnionFS  FS: File System 联合文件系统
分层文件系统
OverLayer FS
Copy On Write : 需要修改原来在 image层中的数据， 会在拷贝一份到最上层，再修改，所以不会影响底层， 如果是删除，会增加一个 白障， 遮挡下层的数据
