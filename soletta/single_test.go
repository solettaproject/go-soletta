package soletta

func process(node FlowNode, port uint16, packet FlowPacket, data interface{}) {
	v, _ := packet.GetBool()
	println("true xor false =", v, "\nPassed data =", *data.(*int))
}

func Example_singleFlow() {
	Init()

	t := NewFlowNodeType("boolean/xor")
	input_ports := []uint16{t.GetInputPort("IN[0]"), t.GetInputPort("IN[1]")}
	output_ports := []uint16{t.GetOutputPort("OUT")}
	data := 11
	n := NewSingleFlow("seconds", *t, input_ports, output_ports, nil, process, &data)

	n.SendPacket("Bool", true, input_ports[0])
	n.SendPacket("Bool", false, input_ports[1])

	Run()

	n.Destroy()

	Shutdown()
}
