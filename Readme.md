# go-soletta #

Provides the go bindings for [Soletta library][]

## Usage ##

```go
import "github.com/kaspersky/go-soletta/soletta"
```

Construct a new Soletta object, then start and stop the soletta engine:

```go
s := soletta.NewSoletta()
ok := s.Start()
ok = s.Stop()
```

[Soletta library]: https://github.com/solettaproject/soletta
