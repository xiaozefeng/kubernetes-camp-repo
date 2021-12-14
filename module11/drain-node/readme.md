## 驱逐pod

1. 创建json 文件, name需要取pod 的名字
```json
{
  "apiVersion": "policy/v1",
  "kind": "Eviction",
  "metadata": {
    "name": "nginx-deployment-pdb-d969f69b5-wnzdr",
    "namespace": "default"
  }
}
```

2. 使用curl 命令调用k8s api
```bash
curl -v  -k  -H  'Content-type: application/json' --key client.key --cert client.crt https://172.16.112.8:6443/api/v1/namespaces/default/pods/nginx-deployment-pdb-d969f69b5-wnzdr/eviction -d @eviction.json
```
pods/nginx-deployment-pdb-d969f69b5-wnzdr/eviction 这个path中的pod的名称也要替换成json文件中的名称
这里的小细节的是 curl 要加 -k 忽略https 

3. curl需要的 client.crt 和 client.key 需要从 ~/.kube/config  中使用 base64 -d 导出到文件

