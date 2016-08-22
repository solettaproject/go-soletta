package soletta

/*
#include <soletta.h>
#include <sol-flow-builder.h>
*/
import "C"
import "runtime"
import "unsafe"
import "errors"

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

//Exports a port
//portIndex is used for cases where portName is actually an array of ports, pass -1 otherwise
func (fb *FlowBuilder) ExportPort(nodeName, portName string, portIndex int, exportedName string, direction int) (err error) {
	cnodename, cportname, cexportedname := C.CString(nodeName), C.CString(portName), C.CString(exportedName)
	defer C.free(unsafe.Pointer(cportname))
	defer C.free(unsafe.Pointer(cnodename))
	defer C.free(unsafe.Pointer(cexportedname))

	err = nil

	switch direction {
	case FlowPortInput:
		r := C.sol_flow_builder_export_port_in(fb.builder, cnodename, cportname, C.int(portIndex), cexportedname)
		if r < 0 {
			err = errors.New("Could not export input port")
		}
	case FlowPortOutput:
		r := C.sol_flow_builder_export_port_out(fb.builder, cnodename, cportname, C.int(portIndex), cexportedname)
		if r < 0 {
			err = errors.New("Could not export output port")
		}
	}

	return
}

//Adds a new flow node named nodeName of type named nodeType.
//A set of options of form (key, value) can be provided.
//Returns true if successful, false otherwise
func (fb *FlowBuilder) AddNodeByTypeName(nodeName, nodeType string, options map[string]string) error {
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

	r := C.sol_flow_builder_add_node_by_type(fb.builder, cname, cnodeType, coptions)
	if r < 0 {
		return errors.New("Error adding node")
	}
	return nil
}

//Adds a new flow node named nodeName of type fnt
//A set of options of form (key, value) can be provided.
func (fb *FlowBuilder) AddNode(nodeName string, fnt *FlowNodeType, options map[string]string) error {
	cname := C.CString(nodeName)
	defer C.free(unsafe.Pointer(cname))

	copts := mapOptionsToFlowOptions(options)
	r := C.sol_flow_builder_add_node(fb.builder, cname, fnt.ctype, copts)
	if r < 0 {
		return errors.New("Error adding node")
	}

	return nil
}

//Add a connection via port names to the connections specification
//of the resulting constructed flow. The connected nodes has to be
//first added using AddNode. portIndex is used in cases where
//input ports are grouped under the same name, pass -1 otherwise.
func (fb *FlowBuilder) Connect(nodeName1, portName1 string, portIndex1 int, nodeName2, portName2 string, portIndex2 int) (err error) {
	cname1, cport1, cname2, cport2 := C.CString(nodeName1), C.CString(portName1), C.CString(nodeName2), C.CString(portName2)
	defer C.free(unsafe.Pointer(cname1))
	defer C.free(unsafe.Pointer(cport1))
	defer C.free(unsafe.Pointer(cname2))
	defer C.free(unsafe.Pointer(cport2))

	err = nil

	r := C.sol_flow_builder_connect(fb.builder, cname1, cport1, C.int(portIndex1), cname2, cport2, C.int(portIndex2))
	if r < 0 {
		err = errors.New("Failed to make connection")
	}
	return
}

//Retrieves the node type of the builder. From the builder's node type
//can be created the root flow node, using CreateNode.
func (fb *FlowBuilder) GetNodeType() FlowNodeType {
	ret := FlowNodeType{}
	ret.ctype = C.sol_flow_builder_get_node_type(fb.builder)
	return ret
}

func (fb *FlowBuilder) init() {
	fb.builder = C.sol_flow_builder_new()
}

func (fb *FlowBuilder) destroy() {
	C.sol_flow_builder_del(fb.builder)
}
