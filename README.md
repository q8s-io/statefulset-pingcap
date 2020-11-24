# StatefulSet

Advanced StatefulSet from PingCAP.

https://github.com/q8s-io/statefulset-pingcap

## Development operator

```go
package main

import (
	"github.com/q8s-io/statefulset-pingcap/client/apis/apps/v1"
)
```

## run advanced statefulset controller locally

Open a new terminal and run controller:

```
kubectl --kubeconfig $cfgfile apply -f manifests/crd.v1.yaml

kubectl --kubeconfig $cfgfile -n kube-system delete ep advanced-statefulset-controller --ignore-not-found
./output/bin/linux/amd64/cmd/controller-manager --kubeconfig $cfgfile -v 4 --leader-elect-resource-name advanced-statefulset-controller --leader-elect-resource-namespace kube-system

```

## deploy a statefulset

```
kubectl apply -f examples/statefulset.yaml
```

## scale out

Note that `--resource-version` is required for CRD objects.

```
RESOURCE_VERSION=$(kubectl get statefulsets.pingcap.com web -ojsonpath='{.metadata.resourceVersion}')
kubectl scale --resource-version=$RESOURCE_VERSION --replicas=4 statefulsets.pingcap.com web
```

## scale in

```
RESOURCE_VERSION=$(kubectl get statefulsets.pingcap.com web -ojsonpath='{.metadata.resourceVersion}')
kubectl scale --resource-version=$RESOURCE_VERSION --replicas=3 statefulsets.pingcap.com web
```

## scale in at arbitrary position

We should set `delete-slots` annotations and decrement `spec.replicas` at the
same time.

```
kubectl apply -f examples/scale-in-statefulset.yaml 
```
