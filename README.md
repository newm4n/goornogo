# goornogo

Enforce coverage test in your golang project pipeline.

## Install

```text
$ go install github.com/newm4n/goornogo
```

## In your pipeline

```text
go test ./... -covermode=count -coverprofile=coverage.out
goornogo -i coverage.out -c 60
```

Params :

- `i` path to coverage report file
- `c` minimum coverage in percentage, 10 = 10%, 45.6 = 45.6% 

If coverage is above minimum coverage, goornogo will exit with code 0.
If bellow the minimum coverage, it will exit with code 1, failing your pipeline.

### Travis-CI

```yaml
language: go

go:
  - 1.13.x

script:
  - go install github.com/newm4n/goornogo 
  - go test ./... -v -covermode=count -coverprofile=coverage.out
  - goornogo -i coverage.out -c 45.3
```

### Circle-CI

```yaml
version: 2
jobs:
  build:
    docker:
      - image: circleci/golang:1.13

    working_directory: /go/src/github.com/my/beautiful-project
    steps:
      - checkout

      # specify any bash command here prefixed with `run: `
      - run: go get -v -t -d ./...
      - run: go install github.com/newm4n/goornogo
      - run: go test ./... -v -covermode=count -coverprofile=coverage.out
      - run: goornogo -i coverage.out -c 45.3
```

### Azure DevOps

```yaml
trigger:
  - master

variables:
  GO111MODULE: 'on'
  GOBIN:  '$(GOPATH)/bin' # Go binaries path
  GOROOT: '/usr/local/go1.13' # Go installation path
  GOPATH: '$(system.defaultWorkingDirectory)/gopath' # Go workspace path
  modulePath: '$(GOPATH)/src/github.com/my/beautiful-project' # Path to the module's code

steps:
  - script: |
      mkdir -p '$(GOBIN)'
      mkdir -p '$(GOPATH)/pkg'
      mkdir -p '$(modulePath)'
      shopt -s extglob
      mv !(gopath) '$(modulePath)'
      echo '##vso[task.prependpath]$(GOBIN)'
      echo '##vso[task.prependpath]$(GOROOT)/bin'
      echo '##vso[task.setvariable variable=path]$(PATH):$(GOBIN)'
    displayName: 'Set up the Go workspace'
  - script: go get -v -t -d ./...
    workingDirectory: '$(modulePath)'
    displayName: 'go get dependencies'
  - script: go build -v ./...
    workingDirectory: '$(modulePath)'
    displayName: 'Build'
  - script: |
      go install github.com/newm4n/goornogo
      go test ./... -v -covermode=count -coverprofile=coverage.out
      goornogo -i coverage.out -c 45.3
    workingDirectory: '$(modulePath)'
    displayName: 'Run tests'
```

### Gitlab-CI

```yaml
image: golang:1.13

stages:
  - build
  - test

build:
  script:
    - go build ./...

test:
  script:
    - go install github.com/newm4n/goornogo
    - go test ./... -v -covermode=count -coverprofile=coverage.out
    - goornogo -i coverage.out -c 45.3
```

### Github Action

```yaml
on:
  pull_request:
    branches:
      - master
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Install Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.13
      - name: Checkout code
        uses: actions/checkout@v2
      - name: Fetching dependencies
        run : go get -v -t -d ./...
      - name: Install Goornogo
        run : go install github.com/newm4n/goornogo
      - name: Execute test
        run : |
          go test ./... -v -covermode=count -coverprofile=coverage.out
          goornogo -i coverage.out -c 45.3
```

## How you do that ?

Goornogo logic is based on 

[This code](https://github.com/golang/go/blob/2bc8d90fa21e9547aeb0f0ae775107dc8e05dc0a/src/cmd/cover/html.go#L96)
and [this code](https://github.com/golang/go/blob/2bc8d90fa21e9547aeb0f0ae775107dc8e05dc0a/src/cmd/cover/profile.go#L56) 
