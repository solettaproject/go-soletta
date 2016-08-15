package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
#include <sol-flow-single.h>

extern void goSingleFlowProcessCallback(void *data, struct sol_flow_node *node, uint16_t port, struct sol_flow_packet *packet);

static struct sol_flow_node *single_flow_bridge(const char *id, const struct sol_flow_node_type *type, const struct sol_flow_node_options *options, const uint16_t *connected_ports_in, const uint16_t *connected_ports_out, void *data)
{
    return sol_flow_single_new(id, type, options, connected_ports_in, connected_ports_out, (void *) goSingleFlowProcessCallback, data);
}
*/
import "C"
import "unsafe"

type SingleFlowNode struct {
	*FlowNode
}

//Creates a flow from a single node
//This is useful for scenaries when usage of a single node is desired,
//manually feeding and processing packets on the node's ports.
func NewSingleFlowNode(nodeName string, nodeType FlowNodeType, inputPorts, outputPorts []uint16, options map[string]string, cb SingleFlowProcessCallback, data interface{}) *SingleFlowNode {
	cname := C.CString(nodeName)
	defer C.free(unsafe.Pointer(cname))

	var coptions *C.struct_sol_flow_node_options

	strvOptions := newstrvOptions(options)
	success := true
	if strvOptions != nil {
		defer strvOptions.destroy()
		namedOptions := C.struct_sol_flow_node_named_options{}
		r := C.sol_flow_node_named_options_init_from_strv(&namedOptions, nodeType.nodeType, strvOptions.cstrvOptions)
		if r == 0 {
			defer C.sol_flow_node_named_options_fini(&namedOptions)
			C.sol_flow_node_options_new(nodeType.nodeType, &namedOptions, &coptions)
			defer C.sol_flow_node_options_del(nodeType.nodeType, coptions)
		} else {
			success = false
		}
	}
	if !success {
		/* Assume this is a Go custom type */
		coptions = mapOptionsToFlowOptions(options)
	}

	/* Create C array parameters for input and output ports */
	step := unsafe.Sizeof((C.uint16_t)(0))
	cinputPorts := C.malloc(C.size_t(uintptr(len(inputPorts)+1) * step))
	coutputPorts := C.malloc(C.size_t(uintptr(len(outputPorts)+1) * step))
	defer C.free(cinputPorts)
	defer C.free(coutputPorts)
	pindexIn, pindexOut := uintptr(cinputPorts), uintptr(coutputPorts)
	for _, inputPort := range inputPorts {
		*(*C.uint16_t)(unsafe.Pointer(pindexIn)) = C.uint16_t(inputPort)
		pindexIn += step
	}
	for _, outputPort := range outputPorts {
		*(*C.uint16_t)(unsafe.Pointer(pindexOut)) = C.uint16_t(outputPort)
		pindexOut += step
	}
	*(*C.uint16_t)(unsafe.Pointer(pindexIn)) = C.UINT16_MAX
	*(*C.uint16_t)(unsafe.Pointer(pindexOut)) = C.UINT16_MAX

	p := mapPointer(&singleFlowPacked{cb, data})
	cnode := C.sol_flow_single_new(cname, nodeType.nodeType, coptions, (*C.uint16_t)(cinputPorts), (*C.uint16_t)(coutputPorts), (*[0]byte)(C.goSingleFlowProcessCallback), unsafe.Pointer(p))

	if cnode == nil {
		return nil
	}

	return &SingleFlowNode{&FlowNode{cnode}}
}

//Connects (enables) the specified port
func (sf *SingleFlowNode) ConnectPort(portIndex uint16, direction int) {
	switch direction {
	case FlowPortInput:
		C.sol_flow_single_connect_port_in(sf.FlowNode.cnode, C.uint16_t(portIndex))
	case FlowPortOutput:
		C.sol_flow_single_connect_port_out(sf.FlowNode.cnode, C.uint16_t(portIndex))
	}
}

//Disconnects (disables) the specified port
func (sf *SingleFlowNode) DisconnectPort(portIndex uint16, direction int) {
	switch direction {
	case FlowPortInput:
		C.sol_flow_single_disconnect_port_in(sf.FlowNode.cnode, C.uint16_t(portIndex))
	case FlowPortOutput:
		C.sol_flow_single_disconnect_port_out(sf.FlowNode.cnode, C.uint16_t(portIndex))
	}
}

//Frees the resources associated with the single flow node
func (sf *SingleFlowNode) Destroy() {
	sf.FlowNode.Destroy()
}

//Callback triggered whenever there is data on node's ports
type SingleFlowProcessCallback func(node FlowNode, port uint16, packet FlowPacket, data interface{})

type singleFlowPacked struct {
	cb   SingleFlowProcessCallback
	data interface{}
}

//export goSingleFlowProcessCallback
func goSingleFlowProcessCallback(data unsafe.Pointer, cnode *C.struct_sol_flow_node, cport C.uint16_t, cpacket *C.struct_sol_flow_packet) {
	p := getPointerMapping(uintptr(data)).(*singleFlowPacked)

	flowNode := FlowNode{cnode}
	port := uint16(cport)
	flowPacket := FlowPacket{cpacket}

	p.cb(flowNode, port, flowPacket, p.data)
}
