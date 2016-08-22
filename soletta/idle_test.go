package soletta_test

import "github.com/solettaproject/go-soletta/soletta"
import "testing"

type idleContext struct {
	data   int
	handle soletta.IdleHandle
}

func idleCb1(data interface{}) bool {
	tc := data.(*idleContext)
	if tc.data == 10 {
		soletta.RemoveIdle(tc.handle)
		return false
	}
	tc.data++
	return true
}

func idleCb2(data interface{}) bool {
	tc := data.(*idleContext)
	if tc.data == 20 {
		soletta.RemoveIdle(tc.handle)
		soletta.Quit()
		return false
	}
	tc.data++
	return true
}

func TestIdle(test *testing.T) {
	var ic1, ic2 idleContext

	soletta.Init()

	ic1.handle = soletta.AddIdle(idleCb1, &ic1)
	ic2.handle = soletta.AddIdle(idleCb2, &ic2)

	soletta.Run()
	soletta.Shutdown()

	if ic1.data != 10 {
		test.Fail()
	}

	if ic2.data != 20 {
		test.Fail()
	}
}
