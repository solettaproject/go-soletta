# go-soletta #

[![Build Status](https://semaphoreci.com/api/v1/solettaproject/go-soletta/branches/master/shields_badge.svg)](https://semaphoreci.com/solettaproject/go-soletta)<br/>

Provides the go bindings for [Soletta™ Project library][1].

**Documentation:** [![GoDoc](https://godoc.org/github.com/solettaproject/go-soletta/soletta?status.svg)](https://godoc.org/github.com/solettaproject/go-soletta/soletta)

## Deployment ##

```
go get github.com/solettaproject/go-soletta/soletta
```

## Usage ##

```go
import "github.com/solettaproject/go-soletta/soletta"
```

A minimal example:

```go
ok := soletta.Init()
if ok {
    soletta.Run()
    soletta.Shutdown()
}
```

[1]: https://github.com/solettaproject/soletta
