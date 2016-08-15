package soletta

/*
#include <soletta.h>
#include <sol-flow-builder.h>
*/
import "C"
import "runtime"
import "unsafe"

//Represents a handle for a flow builder
type FlowBuilder struct {
	builder *C.struct_sol_flow_builder
}

//Returns a newly constructed and initialized flow builder.
//The builder will be automatically destroyed when no longer needed.
func NewFlowBuilder() *FlowBuilder {
	builder := &FlowBuilder{}
	builder.init()
	runtime.SetFinalizer(builder, func(fb *FlowBuilder) { fb.destroy() })
	return builder
}

//Adds a new flow node named nodeName of type named nodeType.
//A set of options of form (key, value) can be provided.
//Returns true if successful, false otherwise
func (fb *FlowBuilder) AddNodeByTypeName(nodeName, nodeType string, options map[string]string) (ok bool) {
	cname, cnodeType := C.CString(nodeName), C.CString(nodeType)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cnodeType))

	/* Create the node options */
	var coptions **C.char
	strvOptions := newstrvOptions(options)
	if strvOptions != nil {
		defer strvOptions.destroy()
		coptions = strvOptions.cstrvOptions
	}

	ok = true
	r := C.sol_flow_builder_add_node_by_type(fb.builder, cname, cnodeType, coptions)
	if r < 0 {
		ok = false
	}
	return
}

//Adds a new flow node named nodeName of type fnt
//A set of options of form (key, value) can be provided.
//Returns true if successful, false otherwise
func (fb *FlowBuilder) AddNode(nodeName string, fnt *FlowNodeType, options map[string]string) (ok bool) {
	cname := C.CString(nodeName)
	defer C.free(unsafe.Pointer(cname))

	ok = true

	copts := mapOptionsToFlowOptions(options)
	r := C.sol_flow_builder_add_node(fb.builder, cname, fnt.nodeType, copts)
	if r < 0 {
		ok = false
	}

	return
}

//Add a connection via port names to the connections specification
//of the resulting constructed flow. The connected nodes has to be
//first added using AddNode.
func (fb *FlowBuilder) Connect(name1, port1, name2, port2 string) {
	cname1, cport1, cname2, cport2 := C.CString(name1), C.CString(port1), C.CString(name2), C.CString(port2)
	defer C.free(unsafe.Pointer(cname1))
	defer C.free(unsafe.Pointer(cport1))
	defer C.free(unsafe.Pointer(cname2))
	defer C.free(unsafe.Pointer(cport2))
	C.sol_flow_builder_connect(fb.builder, cname1, cport1, -1, cname2, cport2, -1)
}

//Retrieves the node type of the builder. From the builder's node type
//can be created the root flow node, using CreateNode.
func (fb *FlowBuilder) GetNodeType() FlowNodeType {
	ret := FlowNodeType{}
	ret.nodeType = C.sol_flow_builder_get_node_type(fb.builder)
	return ret
}

func (fb *FlowBuilder) init() {
	fb.builder = C.sol_flow_builder_new()
}

func (fb *FlowBuilder) destroy() {
	C.sol_flow_builder_del(fb.builder)
}
