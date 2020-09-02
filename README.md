# goornogo

Enforce coverage test in your pipeline.

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