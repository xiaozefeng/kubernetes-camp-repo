## install helm

## add repo
helm repo add apollo https://charts.apolloconfig.com
helm search repo apollo

## install apollo config
helm install apollo-service-pro -f config-service-values.yaml -n sre apollo/apollo-service 

#### output
```text
NAME: apollo-service-pro
LAST DEPLOYED: Sun Jan  2 15:06:28 2022
NAMESPACE: sre
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
Meta service url for current release:
  echo http://apollo-service-pro-apollo-configservice.sre:8080

For local test use:
  export POD_NAME=$(kubectl get pods --namespace sre -l "app=apollo-service-pro-apollo-configservice" -o jsonpath="{.items[0].metadata.name}")
  echo http://127.0.0.1:8080
  kubectl --namespace sre port-forward $POD_NAME 8080:8080

Urls registered to meta service:
Config service: http://apollo-service-pro-apollo-configservice.sre:8080
Admin service: http://apollo-service-pro-apollo-adminservice.sre:8090
```

## install apollo portal
helm install apollo-portal -f portal-service-values.yaml -n sre apollo/apollo-portal 