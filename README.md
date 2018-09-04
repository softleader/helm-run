[![Build Status](https://travis-ci.org/softleader/helm-run.svg?branch=master)](https://travis-ci.org/softleader/helm-run)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/softleader/helm-runs/blob/master/LICENSE)
[![Build Status](https://github-basic-badges.herokuapp.com/release/softleader/helm-run.svg)](https://github.com/softleader/helm-run/releases)

# Helm Run Plugin

Run command in container-based environment, commands store on [softleader/dockerfile/helm](https://github.com/softleader/dockerfile/tree/master/helm
)

![](./artitecture.svg)

> helm-run requires docker installed

## Install

Fetch the latest binary release of helm-run and install it:
 
```sh
$ helm plugin install https://github.com/softleader/helm-run
```

### Install FAQ

This section tracks some of the more frequently encountered issues.

#### Q: A required privilege is not held by the client

權限不夠, 請以系統管理員身份 (windows) 或 sudo (linux) 安裝

#### Q: The system cannot find the file specified

安裝 helm 後尚未 initial, 請執行 `helm init -c` 後再次安裝

#### Q: exec: "sh": executable file not found in %PATH%

Windows 使用者必須開著 docker 進行 plugin 的安裝

## Usage

```sh
$ helm run [flags] COMMAND [ARGS]
```

### Flags

```sh
Flags:
      --always-pull-image   always pull image before running command
      --dos2unix            convert FILE from DOS to Unix format (default true)
  -h, --help                help for helm
  -l, --local               command store on local storage, not on github
  -m, --make                executed via GNU make utility, not bash
      --owner string        github owner of command (default "softleader")
      --path-base string    github path base of command (default "helm")
      --repo string         github repository of command (default "dockerfile")
      --rm                  automatically remove the container when it exits (default true)
      --token string        github access token of command for private repositories
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
