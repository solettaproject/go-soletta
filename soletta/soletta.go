package soletta

/*
#cgo CFLAGS: -I/usr/include/soletta/
#cgo LDFLAGS: -lsoletta
#include "soletta.h"
*/
import "C"

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
