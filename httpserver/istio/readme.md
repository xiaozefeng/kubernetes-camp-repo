## 以 istio gateway 部署 httpserver https服务

### 部署服务
```shell
k create ns sec
k label ns sec istio-injection=enabled
k apply -f httpsserver.yaml -n sec
```

### 生成证书
```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
k apply -n istio-system secret tls https-server-credential --key=cncamp.io.key --cert=cncamp.io.crt
k apply 0f istio-specs.yaml -n sec
```

### 获取 istio-ingressgateway 这个 service 的 service ip
```shell
k get svc -n istio-system |grep istio-ingressgateway
```

### 通过https访问
```shell
curl --resolve httpsserver.cncamp.io:443:$INGRESS_IP https://httpsserver.cacamp.io/httpsrever/healthz -v -k
```



