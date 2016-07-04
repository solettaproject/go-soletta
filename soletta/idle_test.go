package soletta

import "testing"

type idleContext struct {
	data   int
	handle IdleHandle
}

func idleCb1(data interface{}) bool {
	tc := data.(*idleContext)
	if tc.data == 10 {
		RemoveIdle(tc.handle)
		return false
	}
	tc.data++
	return true
}

func idleCb2(data interface{}) bool {
	tc := data.(*idleContext)
	if tc.data == 20 {
		RemoveIdle(tc.handle)
		Quit()
		return false
	}
	tc.data++
	return true
}

func TestIdle(t *testing.T) {
	var ic1, ic2 idleContext

	Init()

	ic1.handle = AddIdle(idleCb1, &ic1)
	ic2.handle = AddIdle(idleCb2, &ic2)

	Run()
	Shutdown()

	if ic1.data != 10 {
		t.Fail()
	}

	if ic2.data != 20 {
		t.Fail()
	}
}
