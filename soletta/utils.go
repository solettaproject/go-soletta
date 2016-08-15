package soletta

/*
#include <soletta.h>
#include <sol-flow.h>

struct simple_type_options
{
    struct sol_flow_node_options base;
    void *coptions;
};

static struct simple_type_options *create_options(void *options)
{
    struct simple_type_options *ret = malloc(sizeof *ret);
    ret->base.api_version = SOL_FLOW_NODE_OPTIONS_API_VERSION;
    ret->coptions = options;
    return ret;
}

static struct sol_flow_node_options *get_base_options(struct simple_type_options *options)
{
    return &options->base;
}
*/
import "C"
import "unsafe"

type strvOptions struct {
	cstrvOptions **C.char
	count        int
}

func newstrvOptions(options map[string]string) *strvOptions {
	if options == nil {
		return nil
	}
	step := unsafe.Sizeof((*C.char)(nil))
	coptions := C.malloc(C.size_t(uintptr(len(options)+1) * step))
	pindex := uintptr(coptions)
	for key, value := range options {
		coption := C.CString(key + "=" + value)
		*(**C.char)(unsafe.Pointer(pindex)) = coption
		pindex += step
	}
	*(**C.char)(unsafe.Pointer(pindex)) = nil

	return &strvOptions{(**C.char)(coptions), len(options)}
}

func (so *strvOptions) destroy() {
	step := unsafe.Sizeof((*C.char)(nil))
	pindex := uintptr(unsafe.Pointer(so.cstrvOptions))
	for i := 0; i < so.count; i++ {
		C.free(unsafe.Pointer(*(**C.char)(unsafe.Pointer(pindex))))
		pindex += step
	}
	C.free(unsafe.Pointer(so.cstrvOptions))
}

func mapOptionsToFlowOptions(options map[string]string) *C.struct_sol_flow_node_options {
	p := mapPointer(options)
	opts := C.create_options(unsafe.Pointer(p))
	return C.get_base_options(opts)
}

func flowOptionsToMapOptions(coptions *C.struct_sol_flow_node_options) map[string]string {
	if coptions == nil {
		return nil
	}
	i := getPointerMapping(uintptr(unsafe.Pointer((*C.struct_simple_type_options)(unsafe.Pointer(coptions)).coptions)))
	if i == nil {
		return nil
	}
	return i.(map[string]string)
}

func getSimpleTypeOptionsSize() uintptr {
	return unsafe.Sizeof(C.struct_simple_type_options{})
}
