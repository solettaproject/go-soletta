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

//Creates a flow from a single node
//This is useful for scenaries when usage of a single node is desired,
//manually feeding and processing packets on the node's ports.
func NewSingleFlow(nodeName string, nodeType FlowNodeType, inputPorts, outputPorts []uint16, options map[string]string, cb SingleFlowProcessCallback, data interface{}) *FlowNode {
	cname := C.CString(nodeName)
	defer C.free(unsafe.Pointer(cname))

	var coptions *C.struct_sol_flow_node_options

	strvOptions := newstrvOptions(options)
	if strvOptions != nil {
		defer strvOptions.destroy()

		namedOptions := C.struct_sol_flow_node_named_options{}
		C.sol_flow_node_named_options_init_from_strv(&namedOptions, nodeType.nodeType, strvOptions.cstrvOptions)
		defer C.sol_flow_node_named_options_fini(&namedOptions)
		C.sol_flow_node_options_new(nodeType.nodeType, &namedOptions, &coptions)
		defer C.sol_flow_node_options_del(nodeType.nodeType, coptions)
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
	flowNode := C.single_flow_bridge(cname, nodeType.nodeType, coptions, (*C.uint16_t)(cinputPorts), (*C.uint16_t)(coutputPorts), unsafe.Pointer(p))

	return &FlowNode{flowNode}
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
