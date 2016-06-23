# go-soletta #

Provides the go bindings for [Soletta library][1].


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
