package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
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

//Frees the resources associated with the flow node
func (fn *FlowNode) Destroy() {
	C.sol_flow_node_del(fn.node)
}

//Creates a flow node of this node type.
func (fnt *FlowNodeType) CreateNode(parent *FlowNode, id string, options FlowOptions) FlowNode {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	var cpnode *C.struct_sol_flow_node
	if parent != nil {
		cpnode = parent.node
	}
	cnode := C.sol_flow_node_new(cpnode, cid, fnt.nodeType, nil)
	return FlowNode{node: cnode}
}

//Frees the resources associated with the flow node type.
func (fnt *FlowNodeType) Destroy() {
	C.sol_flow_node_type_del(fnt.nodeType)
}
