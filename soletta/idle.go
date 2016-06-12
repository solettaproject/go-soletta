package soletta

/*
#cgo CFLAGS: -I/usr/include/soletta/
#cgo LDFLAGS: -lsoletta
#include "soletta.h"

extern bool goIdle(void *data);

static struct sol_idle *idle_bridge(void *data)
{
    return sol_idle_add(goIdle, data);
}
*/
import "C"
import "unsafe"

type IdleCallback func(i interface{}) bool

func AddIdle(cb IdleCallback, i interface{}) interface{} {
	p := unsafe.Pointer(&idlePacked{cb, i})
	return C.idle_bridge(p)
}

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
