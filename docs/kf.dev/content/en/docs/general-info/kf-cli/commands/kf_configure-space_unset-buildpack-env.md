---
title: "kf configure-space unset-buildpack-env"
slug: kf-configure-space-unset-buildpack-env
url: /docs/general-info/kf-cli/commands/kf-configure-space-unset-buildpack-env/
---
## kf configure-space unset-buildpack-env

Unset an environment variable for buildpack builds in a space.

### Synopsis

Unset an environment variable for buildpack builds in a space.

```
kf configure-space unset-buildpack-env SPACE_NAME ENV_VAR_NAME [flags]
```

### Examples

```
  kf configure-space unset-buildpack-env my-space JDK_VERSION
```

### Options

```
  -h, --help   help for unset-buildpack-env
```

### Options inherited from parent commands

```
      --config string       Config file (default is $HOME/.kf)
      --kubeconfig string   Kubectl config file (default is $HOME/.kube/config)
      --namespace string    Kubernetes namespace to target
```

### SEE ALSO

* [kf configure-space](/docs/general-info/kf-cli/commands/kf-configure-space/)	 - Set configuration for a space

