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

//Adds a function to be called when the application goes idle.
//cb is called with the context argument when the main loop reaches the idle state.
//A return value of false will get the idler removed.
//Returns a handler which can be used to delete the idler at a later point.
func AddIdle(cb IdleCallback, context interface{}) interface{} {
	p := unsafe.Pointer(&idlePacked{cb, context})
	return C.idle_bridge(p)
}

//Deletes a previously registered idler, based on its handle.
func RemoveIdle(handle interface{}) {
	C.sol_idle_del(handle.(*C.struct_sol_idle))
}

type idlePacked struct {
	cb   IdleCallback
	data interface{}
}

//export goIdle
func goIdle(data unsafe.Pointer) C.bool {
	p := (*idlePacked)(data)
	return C.bool(p.cb(p.data))
}
