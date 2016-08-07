package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
#include <sol-flow-resolver.h>
*/
import "C"
import "unsafe"
import "image/color"
import "time"

//Represents a node in the flow based programming paradigm
type FlowNode struct {
	cnode *C.struct_sol_flow_node
}

//A node type carries information about node operations,
//input/output ports, options, descriptions etc.
//Represents a blueprint for constructing flow nodes.
type FlowNodeType struct {
	nodeType *C.struct_sol_flow_node_type
}

//Represents a port used in flow node connections
type FlowPort struct {
}

//Sends a packet on this port
func (fn *FlowNode) SendPacket(packetType string, value interface{}, port uint16) {
	cport := C.uint16_t(port)
	switch packetType {
	case "Bool":
		v := C.uchar(0)
		if value.(bool) {
			v = 1
		}
		C.sol_flow_send_bool_packet(fn.cnode, cport, v)
	case "String":
		cstring := C.CString(value.(string))
		defer C.free(unsafe.Pointer(cstring))
		C.sol_flow_send_string_packet(fn.cnode, cport, cstring)
	case "Integer":
		C.sol_flow_send_irange_value_packet(fn.cnode, cport, C.int32_t(value.(int32)))
	case "Double":
		C.sol_flow_send_drange_value_packet(fn.cnode, cport, C.double(value.(float64)))
	case "Byte":
		C.sol_flow_send_byte_packet(fn.cnode, cport, C.uchar(value.(byte)))
	case "Direction":
		dv := value.(DirectionVector)
		cdv := C.struct_sol_direction_vector{C.double(dv.Max), C.double(dv.Min), C.double(dv.X), C.double(dv.Y), C.double(dv.Z)}
		C.sol_flow_send_direction_vector_packet(fn.cnode, cport, &cdv)
	case "Color":
		r, g, b, _ := value.(color.Color).RGBA()
		C.sol_flow_send_rgb_components_packet(fn.cnode, cport, C.uint32_t(r), C.uint32_t(g), C.uint32_t(b))
	case "Location":
		loc := value.(Location)
		C.sol_flow_send_location_components_packet(fn.cnode, cport, C.double(loc.Latitude), C.double(loc.Longitude), C.double(loc.Altitude))
	case "Time":
		t := value.(time.Time)
		ctimespec := C.struct_timespec{C.__time_t(t.Unix()), C.__syscall_slong_t(t.UnixNano())}
		C.sol_flow_send_timestamp_packet(fn.cnode, cport, &ctimespec)
	}
}

//Frees the resources associated with the flow node
func (fn *FlowNode) Destroy() {
	C.sol_flow_node_del(fn.cnode)
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
		cpnode = parent.cnode
	}

	var coptions *C.struct_sol_flow_node_options
	strvOptions := newstrvOptions(options)
	if strvOptions != nil {
		defer strvOptions.destroy()
		namedOptions := C.struct_sol_flow_node_named_options{}
		C.sol_flow_node_named_options_init_from_strv(&namedOptions, fnt.nodeType, strvOptions.cstrvOptions)
		defer C.sol_flow_node_named_options_fini(&namedOptions)
		C.sol_flow_node_options_new(fnt.nodeType, &namedOptions, &coptions)
		defer C.sol_flow_node_options_del(fnt.nodeType, coptions)
	}

	cnode := C.sol_flow_node_new(cpnode, cid, fnt.nodeType, coptions)
	return &FlowNode{cnode: cnode}
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
