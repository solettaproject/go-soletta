package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
#include <sol-flow-simple-c-type.h>

struct CPortDescription
{
    char *name;
    sol_flow_packet_type *packetType;
    int portType;
};

extern int goProcessEvent(struct sol_flow_node *node, struct sol_flow_simple_c_type_event *ev, void *data);

//TODO change function signature to allow more than 1 input port and 1 output port
static struct sol_flow_node_type *sol_flow_simple_c_type_new_full_wrapper(const char *name, size_t context_data_size, struct CPortDescription *input_port, struct CPortDescription *output_port)
{
    struct sol_flow_node_type *ret = sol_flow_simple_c_type_new_full(name, context_data_size, sizeof(struct sol_flow_node_options), (void *) goProcessEvent, input_port->name, input_port->packetType, input_port->portType, output_port->name, output_port->packetType, output_port->portType, NULL);
    return ret;
}

static int get_event_type(const struct sol_flow_simple_c_type_event *ev)
{
    return ev->type;
}
*/
import "C"
import "unsafe"

var mapTypeNameToProcessCallback map[string]ProcessSimpleEventCallback = make(map[string]ProcessSimpleEventCallback)

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

	return ret
}

//Creates a custom node type
func NewSimpleNodeType(name string, ports []PortDescription, cb ProcessSimpleEventCallback) *FlowNodeType {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	packetTypes := map[string]*C.struct_sol_flow_packet_type{"Bool": C.SOL_FLOW_PACKET_TYPE_BOOL}

	inputPort := C.struct_CPortDescription{name: C.CString(ports[0].Name), packetType: packetTypes[ports[0].PacketType], portType: C.SOL_FLOW_SIMPLE_C_TYPE_PORT_TYPE_IN}
	defer C.free(unsafe.Pointer(inputPort.name))
	outputPort := C.struct_CPortDescription{name: C.CString(ports[1].Name), packetType: packetTypes[ports[1].PacketType], portType: C.SOL_FLOW_SIMPLE_C_TYPE_PORT_TYPE_OUT}
	defer C.free(unsafe.Pointer(outputPort.name))

	ctype := C.sol_flow_simple_c_type_new_full_wrapper(cname, 0, &inputPort, &outputPort)
	mapTypeNameToProcessCallback[name] = cb
	return &FlowNodeType{ctype}
}

//Callback for processing events associated with the node
//Return true if no error, false otherwise
type ProcessSimpleEventCallback func(node *FlowNode, event *SimpleFlowEvent, data interface{}) bool

//export goProcessEvent
func goProcessEvent(cnode *C.struct_sol_flow_node, ev *C.struct_sol_flow_simple_c_type_event, data unsafe.Pointer) C.int {
	t := C.sol_flow_node_get_type(cnode)
	name := C.GoString(t.description.name)
	cb := mapTypeNameToProcessCallback[name]
	ret := C.int(0)
	r := cb(&FlowNode{cnode}, newSimpleFlowEvent(ev), nil)

	//Convert the callback return value to boolean
	if !r {
		ret = -1
	}
	return ret
}
