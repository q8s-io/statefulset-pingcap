# StatefulSet

Advanced StatefulSet from PingCAP.

https://github.com/pingcap/advanced-statefulset

## Development operator

```go
package main

import (
	"github.com/q8s-io/statefulset-pingcap/client/apis/apps/v1"
)
```

## Verify

```shell script
make verify
```

## Deploy a statefulset

```shell script
kubectl apply -f examples/statefulset.yaml
```

### scale out

Note that `--resource-version` is required for CRD objects.

```shell script
RESOURCE_VERSION=$(kubectl get statefulsets.pingcap.com web -ojsonpath='{.metadata.resourceVersion}')
kubectl scale --resource-version=$RESOURCE_VERSION --replicas=4 statefulsets.pingcap.com web
```

### scale in

```shell script
RESOURCE_VERSION=$(kubectl get statefulsets.pingcap.com web -ojsonpath='{.metadata.resourceVersion}')
kubectl scale --resource-version=$RESOURCE_VERSION --replicas=3 statefulsets.pingcap.com web
```

### scale in at arbitrary position

We should set `delete-slots` annotations and decrement `spec.replicas` at the
same time.

```shell script
kubectl apply -f examples/scale-in-statefulset.yaml 
```
