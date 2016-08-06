package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
#include <sol-flow-resolver.h>
*/
import "C"
import "unsafe"

//Represents a node in the flow based programming paradigm
type FlowNode struct {
	node *C.struct_sol_flow_node
}

//A node type carries information about node operations,
//input/output ports, options, descriptions etc.
//Represents a blueprint for constructing flow nodes.
type FlowNodeType struct {
	nodeType *C.struct_sol_flow_node_type
}

//Represents a collection of options for flow configuration
type FlowOptions struct {
}

//Represents a port used in flow node connections
type FlowPort struct {
}

//Frees the resources associated with the flow node
func (fn *FlowNode) Destroy() {
	C.sol_flow_node_del(fn.node)
}

//Creates a new node type by name
func NewFlowNodeType(typeName string) *FlowNodeType {
	ret := &FlowNodeType{}

	cname := C.CString(typeName)
	defer C.free(unsafe.Pointer(cname))

	namedOptions := C.struct_sol_flow_node_named_options{}
	C.sol_flow_resolve(C.sol_flow_get_builtins_resolver(), cname, &ret.nodeType, &namedOptions)
	if ret.nodeType == nil {
		C.sol_flow_resolve(nil, cname, &ret.nodeType, &namedOptions)
	}

	if ret.nodeType == nil {
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
		cpnode = parent.node
	}

	strvOptions := newstrvOptions(options)
	defer strvOptions.destroy()
	namedOptions := C.struct_sol_flow_node_named_options{}
	C.sol_flow_node_named_options_init_from_strv(&namedOptions, fnt.nodeType, strvOptions.cstrvOptions)
	defer C.sol_flow_node_named_options_fini(&namedOptions)
	var coptions *C.struct_sol_flow_node_options
	C.sol_flow_node_options_new(fnt.nodeType, &namedOptions, &coptions)
	defer C.sol_flow_node_options_del(fnt.nodeType, coptions)

	cnode := C.sol_flow_node_new(cpnode, cid, fnt.nodeType, coptions)
	return &FlowNode{node: cnode}
}

//Gets an input port by name
func (fnt *FlowNodeType) GetInputPort(name string) uint16 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return uint16(C.sol_flow_node_find_port_in(fnt.nodeType, cname))
}

//Gets an output port by name
func (fnt *FlowNodeType) GetOutputPort(name string) uint16 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return uint16(C.sol_flow_node_find_port_out(fnt.nodeType, cname))
}

//Frees the resources associated with the flow node type.
func (fnt *FlowNodeType) Destroy() {
	C.sol_flow_node_type_del(fnt.nodeType)
}
