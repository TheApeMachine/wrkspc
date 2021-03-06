## wrkspc

A dynamically building workspace for people like me, not you.

### Synopsis


Workspace uses container and cluster technology to dynamically build your tried and true
toolset around you on any machine that contains the binary. Nothing else should
be needed to install.


### Options

```
      --config string   config file (default is $HOME/.wrkspc.yml) (default ".wrkspc.yml")
  -h, --help            help for wrkspc
  -k, --kube            Run in Kubernetes cluster if true (will create one if none exists).
      --viper           use Viper for configuration (default true)
```

### SEE ALSO

* [wrkspc amsh](wrkspc_amsh.md)	 - Ape Machine Shell is an interactive console environment that interfaces with wrkspc.
* [wrkspc completion](wrkspc_completion.md)	 - generate the autocompletion script for the specified shell
* [wrkspc run](wrkspc_run.md)	 - Proxies a command through wrkspc so it will download and run the relevant container.
* [wrkspc serve](wrkspc_serve.md)	 - Serve wrkspc as a service.

###### Auto generated by spf13/cobra on 28-Nov-2021
