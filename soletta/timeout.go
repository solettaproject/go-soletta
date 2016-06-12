package soletta

/*
#cgo CFLAGS: -I/usr/include/soletta/
#cgo LDFLAGS: -lsoletta
#include "soletta.h"

extern bool goTimeout(void *data);

static struct sol_timeout *timeout_bridge(void *data, int ms)
{
    return sol_timeout_add(ms, goTimeout, data);
}
*/
import "C"
import "unsafe"

type TimeoutCallback func(data interface{}) bool

func AddTimeout(cb TimeoutCallback, data interface{}, timeout int) interface{} {
	return C.timeout_bridge(unsafe.Pointer(&timeoutPacked{cb, data}), C.int(timeout))
}

func RemoveTimeout(handle interface{}) {
	C.sol_timeout_del(handle.(*C.struct_sol_timeout))
}

type timeoutPacked struct {
	cb   TimeoutCallback
	data interface{}
}

//export goTimeout
func goTimeout(data unsafe.Pointer) C.bool {
	p := (*timeoutPacked)(data)
	return C.bool(p.cb(p.data))
}
