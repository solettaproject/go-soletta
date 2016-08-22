package soletta

/*
#include <soletta.h>
#include <sol-flow.h>

#include "simple_type_c.h"

extern int goProcessEvent(struct sol_flow_node *node, struct sol_flow_simple_c_type_event *ev, void *data);

static int get_event_type(const struct sol_flow_simple_c_type_event *ev)
{
    return ev->type;
}
*/
import "C"
import "unsafe"

var mapTypeNameToProcessCallback map[string]simplePacked = make(map[string]simplePacked)

const (
	SimpleEventOpen                 int = iota
	SimpleEventClose                int = iota
	SimpleEventConnectInputPort     int = iota
	SimpleEventDisconnectInputPort  int = iota
	SimpleEventProcessInputPort     int = iota
	SimpleEventConnectOutputPort    int = iota
	SimpleEventDisconnectOutputPort int = iota
)

//Represents an event
type SimpleFlowEvent struct {
	cevent       *C.struct_sol_flow_simple_c_type_event
	Type         int
	Port         uint16
	PortName     string
	ConnectionId uint16
	Packet       *FlowPacket
	Options      map[string]string
	Data         interface{}
}

func newSimpleFlowEvent(cevent *C.struct_sol_flow_simple_c_type_event) *SimpleFlowEvent {
	ret := &SimpleFlowEvent{cevent: cevent}
	switch C.get_event_type(cevent) {
	case C.SOL_FLOW_SIMPLE_C_TYPE_EVENT_TYPE_OPEN:
		ret.Type = SimpleEventOpen
	case C.SOL_FLOW_SIMPLE_C_TYPE_EVENT_TYPE_CLOSE:
		ret.Type = SimpleEventClose
	case C.SOL_FLOW_SIMPLE_C_TYPE_EVENT_TYPE_CONNECT_PORT_IN:
		ret.Type = SimpleEventConnectInputPort
	case C.SOL_FLOW_SIMPLE_C_TYPE_EVENT_TYPE_DISCONNECT_PORT_IN:
		ret.Type = SimpleEventDisconnectInputPort
	case C.SOL_FLOW_SIMPLE_C_TYPE_EVENT_TYPE_PROCESS_PORT_IN:
		ret.Type = SimpleEventProcessInputPort
	case C.SOL_FLOW_SIMPLE_C_TYPE_EVENT_TYPE_CONNECT_PORT_OUT:
		ret.Type = SimpleEventConnectOutputPort
	case C.SOL_FLOW_SIMPLE_C_TYPE_EVENT_TYPE_DISCONNECT_PORT_OUT:
		ret.Type = SimpleEventDisconnectOutputPort
	}

	ret.PortName = C.GoString(cevent.port_name)
	ret.Port = uint16(cevent.port)
	ret.ConnectionId = uint16(cevent.conn_id)
	ret.Packet = &FlowPacket{cevent.packet}
	ret.Options = flowOptionsToMapOptions(cevent.options)
	ret.Data = nil

	return ret
}

//Creates a custom node type
func NewSimpleNodeType(name string, ports []PortDescription, cb ProcessSimpleEventCallback) *FlowNodeType {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	packetTypes := map[string]*C.struct_sol_flow_packet_type{"Bool": C.SOL_FLOW_PACKET_TYPE_BOOL, "Integer": C.SOL_FLOW_PACKET_TYPE_IRANGE, "String": C.SOL_FLOW_PACKET_TYPE_STRING}

	/* Create the port array */
	step := unsafe.Sizeof(C.struct_CPortDescription{})
	cports := C.malloc(C.size_t(uintptr(len(ports)) * step))
	defer C.free(unsafe.Pointer(cports))
	for i, port := range ports {
		cport := (*C.struct_CPortDescription)(unsafe.Pointer((uintptr(cports) + uintptr(i)*step)))
		cport.name = C.CString(port.Name)
		defer C.free(unsafe.Pointer(cport.name))
		cport.packet_type = packetTypes[port.PacketType]
		switch port.PortType {
		case FlowPortInput:
			cport.direction = C.SOL_FLOW_SIMPLE_C_TYPE_PORT_TYPE_IN
		case FlowPortOutput:
			cport.direction = C.SOL_FLOW_SIMPLE_C_TYPE_PORT_TYPE_OUT
		}
	}

	ctype := C.sol_flow_simple_c_type_new_full(cname, 0, C.uint16_t(getSimpleTypeOptionsSize()), (*[0]byte)(C.goProcessEvent), (*C.struct_CPortDescription)(cports), C.int(len(ports)))
	if ctype == nil {
		return nil
	}

	mapTypeNameToProcessCallback[name] = simplePacked{cb}
	return &FlowNodeType{ctype}
}

//Encapsulates the callback with the associated data
type simplePacked struct {
	cb ProcessSimpleEventCallback
}

//Callback for processing events associated with the node
//Return true if no error, false otherwise
type ProcessSimpleEventCallback func(node *FlowNode, event *SimpleFlowEvent) bool

//export goProcessEvent
func goProcessEvent(cnode *C.struct_sol_flow_node, ev *C.struct_sol_flow_simple_c_type_event, data unsafe.Pointer) C.int {
	t := C.sol_flow_node_get_type(cnode)
	name := C.GoString(t.description.name)
	packed := mapTypeNameToProcessCallback[name]
	ret := C.int(0)
	r := packed.cb(&FlowNode{cnode}, newSimpleFlowEvent(ev))

	//Convert the callback return value to boolean
	if !r {
		ret = -1
	}
	return ret
}
