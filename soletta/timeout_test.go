package soletta_test

import "github.com/solettaproject/go-soletta/soletta"
import "testing"

type timeoutContext struct {
	data   int
	handle soletta.TimeoutHandle
}

func timeoutCb1(data interface{}) bool {
	tc := data.(*timeoutContext)
	if tc.data == 10 {
		soletta.RemoveTimeout(tc.handle)
		return false
	}
	tc.data++
	return true
}

func timeoutCb2(data interface{}) bool {
	tc := data.(*timeoutContext)
	if tc.data == 20 {
		soletta.RemoveTimeout(tc.handle)
		soletta.Quit()
		return false
	}
	tc.data++
	return true
}

func TestTimeout(test *testing.T) {
	var tc1, tc2 timeoutContext

	soletta.Init()

	tc1.handle = soletta.AddTimeout(timeoutCb1, &tc1, 1)
	tc2.handle = soletta.AddTimeout(timeoutCb2, &tc2, 1)

	soletta.Run()
	soletta.Shutdown()

	if tc1.data != 10 {
		test.Fail()
	}

	if tc2.data != 20 {
		test.Fail()
	}
}
