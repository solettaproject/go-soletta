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
    void *data;
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
    msi->data = data;

    return msi;
}
*/
import "C"
import "unsafe"

type MainloopSource interface {
	GetMainloopSourceAPIVersion() uint16
	Check(interface{}) bool
	Dispatch(interface{})
	Dispose(interface{})
	Prepare(interface{}) bool
}

func AddSource(source MainloopSource, data interface{}) interface{} {
	ret := C.source_bridge(unsafe.Pointer(&sourcePacked{source, data}), C.uint16_t(source.GetMainloopSourceAPIVersion()))
	gMainloopSources[ret.handle] = ret
	return ret.handle
}

func RemoveSource(handle interface{}) {
	h, ok := gMainloopSources[handle.(*C.struct_sol_mainloop_source)]
	if !ok {
		return
	}
	C.sol_mainloop_del_source(handle.(*C.struct_sol_mainloop_source))
	C.free(unsafe.Pointer(h.source))
	C.free(unsafe.Pointer(h))
	delete(gMainloopSources, handle.(*C.struct_sol_mainloop_source))
}

func GetSourceData(handle interface{}) interface{} {
	return (*sourcePacked)(gMainloopSources[handle.(*C.struct_sol_mainloop_source)].data).data
}

var gMainloopSources map[*C.struct_sol_mainloop_source](*C.struct_mainloop_source_info) = make(map[*C.struct_sol_mainloop_source](*C.struct_mainloop_source_info))

type sourcePacked struct {
	source MainloopSource
	data   interface{}
}

//export goCheck
func goCheck(data unsafe.Pointer) C.bool {
	p := (*sourcePacked)(data)
	return C.bool(p.source.Check(p.data))
}

//export goDispatch
func goDispatch(data unsafe.Pointer) {
	p := (*sourcePacked)(data)
	p.source.Dispatch(p.data)
}

//export goDispose
func goDispose(data unsafe.Pointer) {
	p := (*sourcePacked)(data)
	p.source.Dispose(p.data)
}

//export goPrepare
func goPrepare(data unsafe.Pointer) C.bool {
	p := (*sourcePacked)(data)
	return C.bool(p.source.Prepare(p.data))
}
