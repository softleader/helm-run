[![Build Status](https://travis-ci.org/softleader/helm-run.svg?branch=master)](https://travis-ci.org/softleader/helm-run)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/softleader/helm-runs/blob/master/LICENSE)

# Helm Run Plugin

Run command in container-based environment, commands store on [softleader/dockerfile/helm](https://github.com/softleader/dockerfile/tree/master/helm
)

## Install

Fetch the latest binary release of helm-run and install it:
 
```sh
$ helm plugin install https://github.com/softleader/helm-run
```

## Usage
 
```sh
$ helm run [flags] COMMAND [ARGS]
```

### Flags

```sh
Flags:
      --always-pull-image   always pull image before running command
  -h, --help                help for helm
  -i, --image string        image for running command (default "softleader/helm")
      --rm                  automatically remove the container when it exits
```
