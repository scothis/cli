---
id: riff-gateway-status
title: "riff gateway status"
---
## riff gateway status

show gateway status

### Synopsis

Display status details for a gateway.

The Ready condition is shown which should include a reason code and a
descriptive message when the status is not "True". The status for the condition
may be: "True", "False" or "Unknown". An "Unknown" status is common while the
gateway rollout is being processed.

```
riff gateway status <name> [flags]
```

### Examples

```
riff streamming gateway status my-gateway
```

### Options

```
  -h, --help             help for status
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff gateway](riff_gateway.md)	 - (experimental) stream gateway

