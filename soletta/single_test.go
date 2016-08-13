package soletta

func singleFlowProcessCallback(node FlowNode, port uint16, packet FlowPacket, data interface{}) {
	v, _ := packet.GetBool()
	println("true xor false =", v, "\nPassed data =", *data.(*int))
}

func Example_singleFlow() {
	Init()

	t := NewFlowNodeType("boolean/xor")
	pi1, _ := t.GetPort("IN[0]", FlowPortInput)
	pi2, _ := t.GetPort("IN[1]", FlowPortInput)
	po1, _ := t.GetPort("OUT", FlowPortOutput)
	inputPorts := []uint16{pi1, pi2}
	outputPorts := []uint16{po1}

	data := 11
	n := NewSingleFlow("seconds", *t, inputPorts, outputPorts, nil, singleFlowProcessCallback, &data)

	n.SendPacket("Bool", true, inputPorts[0])
	n.SendPacket("Bool", false, inputPorts[1])

	Run()

	n.Destroy()

	Shutdown()
}
