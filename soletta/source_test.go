package soletta

import "testing"

type Source struct {
	data   int
	handle MainloopSourceHandle
}

func (t *Source) GetMainloopSourceAPIVersion() uint16 {
	return MainloopSourceAPIVersion
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
		RemoveSource(t.handle)
		Quit()
		return false
	}
	return true
}

func TestSource(t *testing.T) {
	var s Source

	Init()

	s.handle = AddSource(&s, nil)

	Run()
	Shutdown()

	if s.data != 10 {
		t.Fail()
	}
}
