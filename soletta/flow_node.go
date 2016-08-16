package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
*/
import "C"
import "unsafe"
import "image/color"
import "time"

//Represents a node in the flow based programming paradigm
type FlowNode struct {
	cnode *C.struct_sol_flow_node
}

//Retrieves the name (id) of this flow node
func (fn *FlowNode) GetName() string {
	return C.GoString(C.sol_flow_node_get_id(fn.cnode))
}

//Retrieves the type associated with the flow node
//Returns a nil value in case of error
func (fn *FlowNode) GetType() *FlowNodeType {
	ctype := C.sol_flow_node_get_type(fn.cnode)
	if ctype == nil {
		return nil
	}
	return &FlowNodeType{ctype}
}

//Retrieves the port index by its name
func (fn *FlowNode) GetPort(name string, direction int) (portIndex uint16, ok bool) {
	t := fn.GetType()
	if t == nil {
		return C.UINT16_MAX, false
	}
	return t.GetPort(name, direction)
}

//Sets the data associated with the flow node
func (fn *FlowNode) SetData(data interface{}) {
	if data == nil {
		return
	}
	nodeData[fn.cnode] = data
}

//Retrieves the data associated with this flow node
func (fn *FlowNode) GetData() interface{} {
	return nodeData[fn.cnode]
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
	delete(nodeData, fn.cnode)
	C.sol_flow_node_del(fn.cnode)
}

var nodeData map[*C.struct_sol_flow_node]interface{} = make(map[*C.struct_sol_flow_node]interface{})
