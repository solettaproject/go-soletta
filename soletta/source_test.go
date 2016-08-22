package soletta_test

import "github.com/solettaproject/go-soletta/soletta"
import "testing"

type Source struct {
	data   int
	handle soletta.MainloopSourceHandle
}

func (t *Source) GetMainloopSourceAPIVersion() uint16 {
	return soletta.MainloopSourceAPIVersion
}

func (t *Source) Check(data interface{}) bool {
	return true
}

func (t *Source) Dispatch(data interface{}) {
	t.data++
}

func (t *Source) Dispose(data interface{}) {
}

func (t *Source) Prepare(data interface{}) bool {
	if t.data == 10 {
		soletta.RemoveSource(t.handle)
		soletta.Quit()
		return false
	}
	return true
}

func TestSource(test *testing.T) {
	var s Source

	soletta.Init()

	s.handle = soletta.AddSource(&s, nil)

	soletta.Run()
	soletta.Shutdown()

	if s.data != 10 {
		test.Fail()
	}
}
