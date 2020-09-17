---
id: riff-pulsar-gateway
title: "riff pulsar-gateway"
---
## riff pulsar-gateway

(experimental) pulsar stream gateway

### Synopsis

The Pulsar gateway encapsulates the address of a streaming gateway and a Pulsar
provisioner instance.

The Pulsar provisioner is responsible for resolving topic addresses in a Pulsar
cluster. The streaming gateway coordinates and standardizes reads and writes to
a Pulsar broker.

### Options

```
  -h, --help   help for pulsar-gateway
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff](riff.md)	 - riff is for functions
* [riff pulsar-gateway create](riff_pulsar-gateway_create.md)	 - create a pulsar gateway of messages
* [riff pulsar-gateway delete](riff_pulsar-gateway_delete.md)	 - delete pulsar gateway(s)
* [riff pulsar-gateway list](riff_pulsar-gateway_list.md)	 - table listing of pulsar gateways
* [riff pulsar-gateway status](riff_pulsar-gateway_status.md)	 - show pulsar gateway status

