package soletta

/*
#include "soletta.h"

extern bool goTimeout(void *data);

static struct sol_timeout *timeout_bridge(void *data, int ms)
{
    return sol_timeout_add(ms, goTimeout, data);
}
*/
import "C"
import "unsafe"

//Describes a timeout callback to be registered with AddTimeout.
type TimeoutCallback func(context interface{}) bool

//Adds a function to be called every timeout milliseconds by the main loop,
//as long as cb returns true. Returns a handler which can be used to delete
//the timeout callback at a later point.
func AddTimeout(cb TimeoutCallback, context interface{}, timeout int) interface{} {
	return C.timeout_bridge(unsafe.Pointer(&timeoutPacked{cb, context}), C.int(timeout))
}

//Deletes a previously registered timeout, based on its handle.
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
