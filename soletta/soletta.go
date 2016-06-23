package soletta

/*
#cgo pkg-config: soletta

#include "soletta.h"
*/
import "C"

var MainloopSourceAPIVersion uint16 = C.SOL_MAINLOOP_SOURCE_TYPE_API_VERSION

func Init() bool {
	r := C.sol_init()
	if r == 0 {
		return true
	}
	return false
}

func Run() bool {
	r := C.sol_run()
	if r == C.EXIT_SUCCESS {
		return true
	}
	return false
}

func Quit() {
	C.sol_quit()
}

func Shutdown() {
	C.sol_shutdown()
}
