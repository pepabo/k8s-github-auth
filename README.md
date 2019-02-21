# k8s-github-auth

webhook srever for kubernetes authentication

## Usage

Running k8s-github-auth with deployment:

```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-github-auth
  labels:
    app: k8s-github-auth
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8s-github-auth
  template:
    metadata:
      labels:
        app: k8s-github-auth
    spec:
      containers:
      - name: k8s-github-auth
        image: rtakaishi/k8s-github-auth
        ports:
        - containerPort: 8443
---
apiVersion: v1
kind: Service
metadata:
  name: k8s-github-auth-dev
  namespace: kube-system
  labels:
    k8s-app: k8s-github-auth
spec:
  selector:
    app: k8s-github-auth
  clusterIP: 10.96.0.20
  ports:
  - name: tcp
    port: 8443
    targetPort: 8443
    protocol: TCP

```

add option to kube-apiserver:

```
--authentication-token-webhook-config-file=/path/to/webhook_config.yaml
```

webhook_config.yaml:

```yaml
---
apiVersion: v1
kind: Config
preferences: {}
clusters:
  - cluster:
      insecure-skip-tls-verify: true
      server: http://10.69.0.20:8443/webhook
    name: webhook
users:
  - name: webhook
contexts:
  - context:
      cluster: webhook
      user: webhook
    name: webhook
current-context: webhook
```

set credential:

```
$ kubectl config set-credentials test --token=tokentokentoken
```

## Update Image

```
$ docker build -t rtakaishi/k8s-github-auth .
```

```
$ docker push rtakaishi/k8s-github-auth
```

