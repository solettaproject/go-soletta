package soletta

/*
#include "soletta.h"

extern bool goCheck(void *data);
extern void goDispatch(void *data);
extern void goDispose(void *data);
extern bool goPrepare(void *data);

struct mainloop_source_info
{
    struct sol_mainloop_source *handle;
    struct sol_mainloop_source_type *source;
};

static struct mainloop_source_info *source_bridge(void *data, uint16_t mainloop_source_api_version)
{
    struct sol_mainloop_source_type *source = malloc(sizeof *source);
    source->api_version = mainloop_source_api_version;
    source->check = goCheck;
    source->dispatch = goDispatch;
    source->dispose = goDispose;
    source->get_next_timeout = NULL;
    source->prepare = goPrepare;

    struct mainloop_source_info *msi = malloc(sizeof *msi);
    msi->handle = sol_mainloop_add_source(source, data);
    msi->source = source;

    return msi;
}
*/
import "C"
import "unsafe"

/*
Interface which represents a mainloop event source.

    GetMainloopSourceAPIVersion
Provides the API version.

    Check
Function to be called to check if there are events to be dispatched.

    Dispatch
Function to be called during main loop iterations if prepare or check returns true.

    Dispose
Function to be called when the source is deleted.

    Prepare
Function to be called to query the next timeout for the next event in this source.
*/
type MainloopSource interface {
	GetMainloopSourceAPIVersion() uint16
	Check(interface{}) bool
	Dispatch(interface{})
	Dispose(interface{})
	Prepare(interface{}) bool
}

//Represents an opaque source handle
type MainloopSourceHandle struct {
	msi  *C.struct_mainloop_source_info
	pack uintptr
}

//Creates a new source of events to the main loop.
func AddSource(source MainloopSource, data interface{}) MainloopSourceHandle {
	p := mapPointer(&sourcePacked{source, data})
	msi := C.source_bridge(unsafe.Pointer(p), C.uint16_t(source.GetMainloopSourceAPIVersion()))
	handle := MainloopSourceHandle{msi, p}
	return handle
}

//Destroy a source of main loop events.
func RemoveSource(handle MainloopSourceHandle) {
	C.sol_mainloop_del_source(handle.msi.handle)
	C.free(unsafe.Pointer(handle.msi.source))
	C.free(unsafe.Pointer(handle.msi))
}

//Retrieve the user data (context) given to the source at creation time.
func GetSourceData(handle MainloopSourceHandle) interface{} {
	p := getPointerMapping(handle.pack).(*sourcePacked)
	return p.data
}

type sourcePacked struct {
	source MainloopSource
	data   interface{}
}

//export goCheck
func goCheck(data unsafe.Pointer) C.bool {
	p := getPointerMapping(uintptr(data)).(*sourcePacked)
	return C.bool(p.source.Check(p.data))
}

//export goDispatch
func goDispatch(data unsafe.Pointer) {
	p := getPointerMapping(uintptr(data)).(*sourcePacked)
	p.source.Dispatch(p.data)
}

//export goDispose
func goDispose(data unsafe.Pointer) {
	p := getPointerMapping(uintptr(data)).(*sourcePacked)
	p.source.Dispose(p.data)
	removePointerMapping(uintptr(data))
}

//export goPrepare
func goPrepare(data unsafe.Pointer) C.bool {
	p := getPointerMapping(uintptr(data)).(*sourcePacked)
	return C.bool(p.source.Prepare(p.data))
}
