# single master installation doc



## precheck
- 我的任意节点操作系统为 CentOS 7.8 或者 CentOS Stream 8
- 我的任意节点 CPU 内核数量大于等于 2，且内存大于等于 4G
- 我的任意节点 hostname 不是 localhost，且不包含下划线、小数点、大写字母
- 我的任意节点都有固定的内网 IP 地址
- 我的任意节点都只有一个网卡，如果有特殊目的，我可以在完成 K8S 安装后再增加新的网卡
- 我的任意节点上 Kubelet使用的 IP 地址 可互通（无需 NAT 映射即可相互访问），且没有防火墙、安全组隔离

### Linux Distributioin Version 
```bash
cat /etc/redhat-release

# output
CentOS Linux release 7.9.2009 (Core)
```

### Hostname
```bash
hostname

# output
uat22

```
** 修改hostname ** 
```bash
# 修改 hostname
hostnamectl set-hostname your-new-host-name
# 查看修改结果
hostnamectl status
# 设置 hostname 解析
echo "127.0.0.1   $(hostname)" >> /etc/hosts
```

### CPU Cores
```bash
lscpu

# output
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                4
```

### RAM 
```bash
free -h 

# output
```

### network check
```bash
ip route show

# output
default via 172.16.112.1 dev eth0 proto dhcp metric 100
169.254.169.0/24 via 172.16.112.1 dev eth0 proto dhcp metric 100
172.16.112.0/20 dev eth0 proto kernel scope link src 172.16.112.11 metric 100

ip addr

# output

1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1450 qdisc mq state UP group default qlen 1000
    link/ether fa:16:3e:10:7d:c1 brd ff:ff:ff:ff:ff:ff
    inet 172.16.112.11/20 brd 172.16.127.255 scope global noprefixroute dynamic eth0
       valid_lft 85066sec preferred_lft 85066sec
    inet6 fe80::d662:97b7:3976:db84/64 scope link noprefixroute
       valid_lft forever preferred_lft forever
```

## Installation
master 和 worker 都要执行
```bash
# 阿里云 docker hub 镜像
export REGISTRY_MIRROR=https://registry.cn-hangzhou.aliyuncs.com
```
### 安装脚本
```bash
#!/bin/bash

# 文档：https://kubernetes.io/docs/setup/production-environment/container-runtimes/#containerd

cat <<EOF | sudo tee /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF


sudo modprobe overlay
sudo modprobe br_netfilter


# Setup required sysctl params, these persist across reboots
cat <<EOF | sudo tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

# Apply sysctl params without reboot
sudo sysctl --system

yum remove -y containerd.io
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install -y containerd.io-1.4.3




sudo mkdir -p /etc/containerd
containerd config default |  tee /etc/containerd/config.toml

sed -i "s#k8s.gcr.io#registry.aliyuncs.com/k8sxio#g"  /etc/containerd/config.toml
sed -i '/containerd.runtimes.runc.options/a\ \ \ \ \ \ \ \ \ \ \ \ SystemdCgroup = true' /etc/containerd/config.toml
sed -i "s#https://registry-1.docker.io#${REGISTRY_MIRROR}#g"  /etc/containerd/config.toml

systemctl daemon-reload
systemctl enable containerd
systemctl restart containerd


yum install -y nfs-utils
yum install -y wget

systemctl stop firewalld
systemctl disable firewalld

setenforce 0
sed -i "s/SELINUX=enforcing/SELINUX=disabled/g" /etc/selinux/config


swapoff -a
yes | cp /etc/fstab /etc/fstab_bak
cat /etc/fstab_bak |grep -v swap > /etc/fstab


at <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
       http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

yum remove -y kubelet kubeadm kubectl

yum install -y kubelet-1.21.4 kubeadm-1.21.4 kubectl-1.21.4

crictl config runtime-endpoint /run/containerd/containerd.sock


systemctl daemon-reload
systemctl enable kubelet && systemctl start kubelet

containerd --version
kubelet --version

```

## Initialation Master Node

```bash
# set master ip
export MASTER_IP =172.16.112.11

# set apiserver DNS name
export APISERVER_NAME=apiserver.uat20.io

# kubernetes pod network subnet 
export POD_SUBNET=10.100.0.0/16

echo "${MASTER_IP}  ${APISERVER_NAME}" >> /etc/hosts

```

### Staring
```bash

# 在脚本出错时执行
set -e 

if [ ${#POD_SUBNET} -eq 0 ] || [ ${#APISERVER_NAME} -eq 0 ]; then
  echo -e "\033[31;1m请确保您已经设置了环境变量 POD_SUBNET 和 APISERVER_NAME \033[0m"
  echo 当前POD_SUBNET=$POD_SUBNET
  echo 当前APISERVER_NAME=$APISERVER_NAME
  exit 1
fi

rm -f ./kubeadm-config.yaml

cat <<EOF > ./kubeadm-config.yaml
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: v${1}
imageRepository: registry.aliyuncs.com/k8sxio
controlPlaneEndpoint: "${APISERVER_NAME}:6443"
networking:
  serviceSubnet: "10.96.0.0/16"
  podSubnet: "${POD_SUBNET}"
  dnsDomain: "cluster.local"
dns:
  type: CoreDNS
  imageRepository: swr.cn-east-2.myhuaweicloud.com${2}
  imageTag: 1.8.0

---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
EOF

echo ""
echo "抓取镜像，请稍候..."
kubeadm config images pull --config=kubeadm-config.yaml
echo ""
echo "初始化 Master 节点"
kubeadm init --config=kubeadm-config.yaml --upload-certs

rm -rf /root/.kube/
mkdir /root/.kube/
cp -i /etc/kubernetes/admin.conf /root/.kube/config

```

### 检查是否成功
```bash
# 等待所有pod处于running状态
watch kubectl get pod -n kube-system -o wide
```


## 安装网络插件
```bash
export POD_SUBNET=10.100.0.0/16
kubectl apply -f https://kuboard.cn/install-script/v1.21.x/calico-operator.yaml
wget https://kuboard.cn/install-script/v1.21.x/calico-custom-resources.yaml
sed -i "s#192.168.0.0/16#${POD_SUBNET}#" calico-custom-resources.yaml
kubectl apply -f calico-custom-resources.yaml
```
