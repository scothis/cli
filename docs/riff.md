---
id: riff
title: "riff"
---
## riff

riff is for functions

### Synopsis

The riff CLI is a client to the projectriff system CRDs. The CRDs
define the riff API.

Before running riff, please install the projectriff system and its dependencies.
See https://projectriff.io/docs/getting-started/

### Options

```
      --config file       config file (default is $HOME/.riff.yaml)
  -h, --help              help for riff
      --kubeconfig file   kubectl config file (default is $HOME/.kube/config)
      --no-color          disable color output in terminals
      --version           display CLI version
```

### SEE ALSO

* [riff completion](riff_completion.md)	 - generate shell completion script
* [riff doctor](riff_doctor.md)	 - check riff's permissions
* [riff gateway](riff_gateway.md)	 - (experimental) stream gateway
* [riff inmemory-gateway](riff_inmemory-gateway.md)	 - (experimental) in-memory stream gateway
* [riff kafka-gateway](riff_kafka-gateway.md)	 - (experimental) kafka stream gateway
* [riff processor](riff_processor.md)	 - (experimental) processors apply functions to messages on streams
* [riff pulsar-gateway](riff_pulsar-gateway.md)	 - (experimental) pulsar stream gateway
* [riff stream](riff_stream.md)	 - (experimental) streams of messages

