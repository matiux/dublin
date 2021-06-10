Dublin
===

Golang porting of PHP [Broadway](https://github.com/broadway/broadway) Library for CQRS + ES

## Project setup
```shell
git clone git@github.com:matiux/dublin.git && cd dublin
go mod tidy
```

## Running test
```shell
go clean -cache
go clean -testcache
go test ./... -v
```