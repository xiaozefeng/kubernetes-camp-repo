## 把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：
1. 如何实现安全保证
    使用 istio的 Gateway 对象的tls能力
    ```yaml
   apiVersion: networking.istio.io/v1beta1
   kind: Gateway
   metadata:
   name: httpsserver
   spec:
   selector:
   istio: ingressgateway
   servers:
    - hosts:
        - httpsserver.cncamp.io
          port:
          name: https-default
          number: 443
          protocol: HTTPS
          tls:
          mode: SIMPLE
          credentialName: https-server-credential
         ```
2. 七层路由规则；
使用 istio的 VirtualService 对象
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: httpsserver
sepc:
  gateways:
    - httpsserver
  hosts:
    - httpsserver.cncamp.io
  http:
    - match:
        - port: 443
          uri: "/httpserver"
      rewrite:
        uri: "/"
      route:
        - destination:
            host: httpserver.default.svc.cluster.local
            port:
              number: 80
```

