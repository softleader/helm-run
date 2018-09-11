[![Build Status](https://travis-ci.org/softleader/helm-run.svg?branch=master)](https://travis-ci.org/softleader/helm-run)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/softleader/helm-runs/blob/master/LICENSE)
[![Build Status](https://github-basic-badges.herokuapp.com/release/softleader/helm-run.svg)](https://github.com/softleader/helm-run/releases)

# Helm Run Plugin

Run bash command in container-based environment.

![](./artitecture.svg)

To run commands that store on [softleader/dockerfile/helm](https://github.com/softleader/dockerfile/tree/master/helm)

	$ helm run package

To run command on local, use `-l`:

	$ helm run -l path/to/script

To execute via GNU make utility, not bash, use `-m`:

	$ helm run -lm path/to/Makefile

> helm-run requires [helm](https://docs.helm.sh/using_helm/#installing-helm) and [docker](https://www.docker.com/get-started) installed

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

Windows 使用者必須開著 docker 進行 plugin 的安裝, 若還是遇到問題, 建議開著 docker 以 bash (如 Git Bash) 執行

#### Q: Error: exec: "C:\\Users\\Default": file does not exist

helm home 預設在使用者目錄下, 但如果使用者名稱有空白字元時會造成 helm 執行錯, 因此建議調整 helm home: 設定環境變數 `HELM_HOME` 指到 `C:\.helm` 後, 重新 initial helm

## Usage

```sh
$ helm run [flags] COMMAND [ARGS]
```

### Flags

```sh
Flags:
      --dos2unix           convert FILE from DOS to Unix format (default true)
  -h, --help               help for helm
  -l, --local              command store on local storage, not on github
  -m, --make               executed via GNU make utility, not bash
      --owner string       github owner of command (default "softleader")
      --path-base string   github path base of command (default "helm")
      --repo string        github repository of command (default "dockerfile")
      --rm                 automatically remove the container when it exits (default true)
      --token string       github access token of command for private repositories
  -U, --update-image       update image before running command
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

### Usage FAQ

#### Q: Error: No such image: softleader/helm

helm run 是 container-based 的 runtime 環境, 預設使用 image 為: [softleader/helm](https://github.com/softleader/dockerfile/tree/master/helm), 第一次執行時請加上 update image 參數: 

```sh
$ helm run -U package
```

> update 之後的 helm run 即不用再下 `-U`

#### Q: /bin/bash: make: command not found

也許是你 local 的 image 新版本過舊但使用到了新的功能, 建議在 helm run 時加上 update image 參數更新 image:

```sh
$ helm run -U package
```

> update 之後的 helm run 即不用再下 `-U`

#### Q: 以上都做了, 但在 windows 還是跑得很奇怪

將 docker 設定中, Network > DNS Server 改成 Fixed 再試試看

![](https://zxtech.files.wordpress.com/2016/09/image3.png)
