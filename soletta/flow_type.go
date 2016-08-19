package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
#include <sol-flow-resolver.h>
*/
import "C"
import "unsafe"

//A node type carries information about node operations,
//input/output ports, options, descriptions etc.
//Represents a blueprint for constructing flow nodes.
type FlowNodeType struct {
	ctype *C.struct_sol_flow_node_type
}

//Creates a new node type by name
func NewFlowNodeType(typeName string) *FlowNodeType {
	ret := &FlowNodeType{}

	cname := C.CString(typeName)
	defer C.free(unsafe.Pointer(cname))

	namedOptions := C.struct_sol_flow_node_named_options{}
	C.sol_flow_resolve(C.sol_flow_get_builtins_resolver(), cname, &ret.ctype, &namedOptions)
	if ret.ctype == nil {
		C.sol_flow_resolve(nil, cname, &ret.ctype, &namedOptions)
	}

	if ret.ctype == nil {
		return nil
	}

	defer C.sol_flow_node_named_options_fini(&namedOptions)

	return ret
}

//Creates a flow node of this node type.
func (fnt *FlowNodeType) CreateNode(parent *FlowNode, id string, options map[string]string) *FlowNode {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	var cpnode *C.struct_sol_flow_node
	if parent != nil {
		cpnode = parent.cnode
	}

	var coptions *C.struct_sol_flow_node_options

	strvOptions := newstrvOptions(options)
	success := true
	if strvOptions != nil {
		defer strvOptions.destroy()
		namedOptions := C.struct_sol_flow_node_named_options{}
		r := C.sol_flow_node_named_options_init_from_strv(&namedOptions, fnt.ctype, strvOptions.cstrvOptions)
		if r == 0 {
			defer C.sol_flow_node_named_options_fini(&namedOptions)
			C.sol_flow_node_options_new(fnt.ctype, &namedOptions, &coptions)
			defer C.sol_flow_node_options_del(fnt.ctype, coptions)
		} else {
			success = false
		}
	}
	if !success {
		/* Assume this is a Go custom type */
		coptions = mapOptionsToFlowOptions(options)
	}

	cnode := C.sol_flow_node_new(cpnode, cid, fnt.ctype, coptions)
	return &FlowNode{cnode: cnode}
}

//Gets an input port by name
func (fnt *FlowNodeType) GetPort(name string, direction int) (portIndex uint16, ok bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	switch direction {
	case FlowPortInput:
		portIndex = uint16(C.sol_flow_node_find_port_in(fnt.ctype, cname))
	case FlowPortOutput:
		portIndex = uint16(C.sol_flow_node_find_port_out(fnt.ctype, cname))
	}

	ok = true
	if portIndex == C.UINT16_MAX {
		ok = false
	}

	return
}

//Retrieves the number of ports
func (fnt *FlowNodeType) GetPortCount(direction int) int {
	switch direction {
	case FlowPortInput:
		return int(fnt.ctype.ports_in_count)
	case FlowPortOutput:
		return int(fnt.ctype.ports_out_count)
	}
	return 0
}

//Frees the resources associated with the flow node type.
func (fnt *FlowNodeType) Destroy() {
	C.sol_flow_node_type_del(fnt.ctype)
}
