package soletta

import "testing"

type timeoutContext struct {
	data   int
	handle TimeoutHandle
}

func timeoutCb1(data interface{}) bool {
	tc := data.(*timeoutContext)
	if tc.data == 10 {
		RemoveTimeout(tc.handle)
		return false
	}
	tc.data++
	return true
}

func timeoutCb2(data interface{}) bool {
	tc := data.(*timeoutContext)
	if tc.data == 20 {
		RemoveTimeout(tc.handle)
		Quit()
		return false
	}
	tc.data++
	return true
}

func TestTimeout(t *testing.T) {
	var tc1, tc2 timeoutContext

	Init()

	tc1.handle = AddTimeout(timeoutCb1, &tc1, 1)
	tc2.handle = AddTimeout(timeoutCb2, &tc2, 1)

	Run()
	Shutdown()

	if tc1.data != 10 {
		t.Fail()
	}

	if tc2.data != 20 {
		t.Fail()
	}
}
