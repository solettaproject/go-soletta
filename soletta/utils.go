package soletta

/*
#include <soletta.h>
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
