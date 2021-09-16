# parzival

[![release-workflow](https://github.com/unfor19/parzival/actions/workflows/release.yml/badge.svg)](https://github.com/unfor19/parzival/actions/workflows/release.yml) [![release](https://img.shields.io/github/v/release/unfor19/parzival?color=green&label=release&logo=go)](https://github.com/unfor19/parzival/releases/latest) [![pre-release](https://img.shields.io/github/v/release/unfor19/parzival?color=orange&include_prereleases&label=pre-release&logo=go)](https://github.com/unfor19/parzival/releases)

**Work In Progress (WIP)**

A CLI that can get/set more than 10 SSM Parameters by path in a single command.

## Getting Started

1. Download the binary file from the releases page, for example [0.0.2](https://github.com/unfor19/parzival/releases/tag/0.0.2)
   - macOS - Intel chips
    ```bash
    PARZIVAL_OS="darwin" && \
    PARZIVAL_ARCH="amd64" && \
    PARZIVAL_VERSION="0.0.2" && \
    curl -sL -o parzival "https://github.com/unfor19/parzival/releases/download/${PARZIVAL_VERSION}/parzival_${PARZIVAL_VERSION}_${PARZIVAL_OS}_${PARZIVAL_ARCH}"
    ```
   - macOS - M1 chips
    ```bash
    PARZIVAL_OS="darwin" && \
    PARZIVAL_ARCH="arm64" && \
    PARZIVAL_VERSION="0.0.2" && \
    curl -sL -o parzival "https://github.com/unfor19/parzival/releases/download/${PARZIVAL_VERSION}/parzival_${PARZIVAL_VERSION}_${PARZIVAL_OS}_${PARZIVAL_ARCH}"
    ```    
   - Linux - amd64
    ```bash
    PARZIVAL_OS="linux" && \
    PARZIVAL_ARCH="amd64" && \
    PARZIVAL_VERSION="0.0.2" && \
    curl -sL -o parzival "https://github.com/unfor19/parzival/releases/download/${PARZIVAL_VERSION}/parzival_${PARZIVAL_VERSION}_${PARZIVAL_OS}_${PARZIVAL_ARCH}"
    ```
   - [Windows WSL2](https://docs.microsoft.com/en-us/windows/wsl/install-win10) - 386
    ```bash
    PARZIVAL_OS="linux" && \
    PARZIVAL_ARCH="386" && \    
    PARZIVAL_VERSION="0.0.2" && \
    curl -sL -o parzival "https://github.com/unfor19/parzival/releases/download/${PARZIVAL_VERSION}/parzival_${PARZIVAL_VERSION}_${PARZIVAL_OS}_${PARZIVAL_ARCH}"
    ```
2. Set permissions to allow execution of `parzival` binary and move to `/usr/local/bin` dir 
   ```bash
   chmod +x parzival && \
   sudo mv parzival "/usr/local/bin/parzival"
   ```
3. Get SSM Parameters by path
   ```bash
   parzival get --region "us-east-1" \
        --output-file-path ".dev_parameters.json" \
        --parameters-path "/myapp/dev/"
   ```
4. Set SSM Parameters according to the output of `Get`
   ```bash
   parzival set --region "us-east-1" \
        --input-file-path ".dev_parameters.json" \
        --parameters-path "/myapp/stg" \
        --prefix-to-replace "/myapp/dev/"
   ```


## Local Development

<details>

<summary>Expand/Collapse</summary>

For local development, we'll use the following services

- [localstack](https://github.com/localstack/localstack) - A fully functional local cloud (AWS) stack

### Requirements

- [Golang 1.16+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) - As part of the test suite, AWS CLI invokes `ssm put-parameter ...`
### Development Process

#### Initial Setup For Golang

1. Place the source code at `$HOME/go/src/github.com/unfor19/parzival`

2. Add the following to `${HOME}/.bash_profile` or `${HOME}/.bashrc`
    ```bash
    export GOPATH=$HOME/go
    export GOROOT=/usr/local/opt/go/libexec
    export PATH=$PATH:$GOPATH/bin:$GOROOT/bin
    ```

#### Run

```
make up-localstack && \
    go run . get --localstack
```

#### Build

```bash
make build
```

#### Test

```bash
make test
```

Skip SSM Parameter creation by setting before running tests

```bash
export SKIP_PARAM_CREATION="true" && \
make test
```

</details>


## Authors

Created and maintained by [Meir Gabay](https://github.com/unfor19)

## License

This project is licensed under the Apache License - see the [LICENSE](https://github.com/unfor19/parzival/blob/master/LICENSE) file for details
