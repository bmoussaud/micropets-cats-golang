# micropets-cats-golang

micropets-cats-golang is a microservice written in `#go-lang` belonging to a larger application [*Micropets Portal*](https://github.com/bmoussaud/micropets-app).

## How to build

* Local using `make` 

```
make clean deps build
build/micropets-cats-golang
```

* Create the image using `Cloud Native Build Packs`

```
make cnb-image
make docker-run
```
* Managed by `Tanzu Application Platform``

```
./apply_workload.sh
curl -k $(k get ksvc cats-golang -o jsonpath='{.status.url}')
```

### How to Configuration

There are 2 levels of configuration
* The main level is help by the `cats-config` secret that points to a JSON File
* The overide of this file using env values syntax offered by the `viper` [library](https://github.com/spf13/viper). The main use case is to be able to have a dedicated secret to set the [Aria Ops for Apps]() token.

The both files are using the `service bindings` to be associate to the microservice, 1st one managed by the workload, 2nd one managed manually to expose the key as environment variables.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cats-config
  labels:
    app.kubernetes.io/part-of: cats
type: Opaque
stringData:
  type: app-configuration
  pets_config.json: |-
    {
      "service": {
        "port": ":8181",
        "listen": "true",
        "mode": "RANDOM_NUMBER",
        "frequencyError": 10,
        "delay": {
          "period": 200,
          "amplitude": 0.3
        },
        "from": "Europ"
      },
      "observability": {
        "enable": true,
        "application": "micropets",
        "service": "cats",
        "cluster": "us-west",
        "shard": "primary",
        "server": "https://binz.wavefront.com",
        "token": "x-y-z-y"
      }
    }
````

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aria-credentials
  labels:
    app.kubernetes.io/name: shared
    app.kubernetes.io/part-of: micropets
type: servicebinding.io/configuration
stringData:
  type: app-configuration-aria
  observability.enable: "true"
  observability.token: x04de315-z31d-4bc0-c123-d8c4d0853dad
```

```yaml
apiVersion: servicebinding.io/v1alpha3
kind: ServiceBinding
metadata:
  name: cats-golang-aria-credentials
  annotations:    
    kapp.k14s.io/change-group: servicebinding.io/ServiceBindings
  labels:
    app.kubernetes.io/name: cats
    app.kubernetes.io/part-of: micropets    
    app.kubernetes.io/component: run
spec:
  name: app-config-aria
  service:
    apiVersion: v1
    kind: Secret
    name: aria-credentials
  workload:
    apiVersion: serving.knative.dev/v1
    kind: Service
    name: cats-golang
  env:
    - key: observability.enable
      name: MP_OBSERVABILITY.ENABLE
    - key: observability.token
      name: MP_OBSERVABILITY.TOKEN

```