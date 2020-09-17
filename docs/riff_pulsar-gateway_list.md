---
id: riff-pulsar-gateway-list
title: "riff pulsar-gateway list"
---
## riff pulsar-gateway list

table listing of pulsar gateways

### Synopsis

List Pulsar gateways in a namespace or across all namespaces.

For detail regarding the status of a single pulsar gateway, run:

    riff streaming pulsar-gateway status <pulsar-gateway-name>

```
riff pulsar-gateway list [flags]
```

### Examples

```
riff streaming pulsar-gateway list
riff streaming pulsar-gateway list --all-namespaces
```

### Options

```
      --all-namespaces   use all kubernetes namespaces
  -h, --help             help for list
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff pulsar-gateway](riff_pulsar-gateway.md)	 - (experimental) pulsar stream gateway

