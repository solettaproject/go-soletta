package soletta

/*
#include "soletta.h"

extern bool goIdle(void *data);

static struct sol_idle *idle_bridge(void *data)
{
    return sol_idle_add(goIdle, data);
}
*/
import "C"
import "unsafe"

//Describes an idle callback to be registered with AddIdle.
type IdleCallback func(context interface{}) bool

//Represents an opaque idler handle
type IdleHandle struct {
	handle *C.struct_sol_idle
	pack   uintptr
}

//Adds a function to be called when the application goes idle.
//cb is called with the context argument when the main loop reaches the idle state.
//A return value of false will get the idler removed.
//Returns a handler which can be used to delete the idler at a later point.
func AddIdle(cb IdleCallback, context interface{}) IdleHandle {
	p := mapPointer(&idlePacked{cb, context})
	return IdleHandle{C.idle_bridge(unsafe.Pointer(p)), p}
}

//Deletes a previously registered idler, based on its handle.
func RemoveIdle(handle IdleHandle) {
	removePointerMapping(handle.pack)
	C.sol_idle_del(handle.handle)
}

type idlePacked struct {
	cb   IdleCallback
	data interface{}
}

//export goIdle
func goIdle(data unsafe.Pointer) C.bool {
	p := getPointerMapping(uintptr(data)).(*idlePacked)
	return C.bool(p.cb(p.data))
}
