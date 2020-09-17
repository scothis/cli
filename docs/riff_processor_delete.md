---
id: riff-processor-delete
title: "riff processor delete"
---
## riff processor delete

delete processor(s)

### Synopsis

Delete one or more processors by name or all processors within a namespace.

The processor will stop processing messages from the input streams and writing
to the output streams. The streams and messages in each stream are preserved.

```
riff processor delete <name(s)> [flags]
```

### Examples

```
riff streaming processor delete my-processor
riff streaming processor delete --all 
```

### Options

```
      --all              delete all processors within the namespace
  -h, --help             help for delete
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --config file       config file (default is $HOME/.riff.yaml)
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
```

### SEE ALSO

* [riff processor](riff_processor.md)	 - (experimental) processors apply functions to messages on streams
