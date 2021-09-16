# parzival

[![release](https://github.com/unfor19/parzival/actions/workflows/release.yml/badge.svg)](https://github.com/unfor19/parzival/actions/workflows/release.yml)

**Work In Progress (WIP)**

A CLI that can get/set more than 10 SSM Parameters by path in a single command.

## Local Development

<details>

<summary>Expand/Collapse</summary>

For local development, we'll use the following services

- [localstack](https://github.com/localstack/localstack) - A fully functional local cloud (AWS) stack
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) - As part of the test suite, AWS CLI invokes `ssm put-parameter ...`

### Requirements

- [Golang 1.16+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

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

</details>


## Authors

Created and maintained by [Meir Gabay](https://github.com/unfor19)

## License

This project is licensed under the Apache License - see the [LICENSE](https://github.com/unfor19/parzival/blob/master/LICENSE) file for details
