---
title: "kf configure-space unset-env"
slug: kf-configure-space-unset-env
url: /docs/general-info/kf-cli/commands/kf-configure-space-unset-env/
---
## kf configure-space unset-env

Unset a space-wide environment variable.

### Synopsis

Unset a space-wide environment variable.

```
kf configure-space unset-env SPACE_NAME ENV_VAR_NAME [flags]
```

### Examples

```
  kf configure-space unset-env my-space ENVIRONMENT
```

### Options

```
  -h, --help   help for unset-env
```

### Options inherited from parent commands

```
      --config string       Config file (default is $HOME/.kf)
      --kubeconfig string   Kubectl config file (default is $HOME/.kube/config)
      --namespace string    Kubernetes namespace to target
```

### SEE ALSO

* [kf configure-space](/docs/general-info/kf-cli/commands/kf-configure-space/)	 - Set configuration for a space

