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
```

### Command

儲存在 [softleader/dockerfile/helm](https://github.com/softleader/dockerfile/tree/master/helm
) 的 shell 檔案名稱, *大小寫是有區分的*

若要執行的 shell 檔案在更下層的目錄中, 則以 `/` 區隔每個目錄, 如在 [softleader/dockerfile/helm](https://github.com/softleader/dockerfile/tree/master/helm
) 的目錄為: 

```sh
.
└── helm
    ├── Dockerfile
    ├── README.md
    ├── package
    └── subdir
        └── mypackage
```

則執行 `mypackage` 的指令將為: 

```sh
$ helm run subdir/mypackage
```

### Args

會被繼續傳入執行 shell 的 args 中
