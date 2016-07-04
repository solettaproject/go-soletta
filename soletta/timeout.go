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

//Represents an opaque timeout handle
type TimeoutHandle struct {
	handle *C.struct_sol_timeout
	pack   uintptr
}

//Adds a function to be called every timeout milliseconds by the main loop,
//as long as cb returns true. Returns a handler which can be used to delete
//the timeout callback at a later point.
func AddTimeout(cb TimeoutCallback, context interface{}, timeout int) TimeoutHandle {
	p := mapPointer(&timeoutPacked{cb, context})
	return TimeoutHandle{C.timeout_bridge(unsafe.Pointer(p), C.int(timeout)), p}
}

//Deletes a previously registered timeout, based on its handle.
func RemoveTimeout(handle TimeoutHandle) {
	removePointerMapping(handle.pack)
	C.sol_timeout_del(handle.handle)
}

type timeoutPacked struct {
	cb   TimeoutCallback
	data interface{}
}

//export goTimeout
func goTimeout(data unsafe.Pointer) C.bool {
	p := getPointerMapping(uintptr(data)).(*timeoutPacked)
	return C.bool(p.cb(p.data))
}
