package soletta_test

import "github.com/solettaproject/go-soletta/soletta"
import "fmt"

func singleFlowProcessCallback(node *soletta.FlowNode, port uint16, packet *soletta.FlowPacket, data interface{}) {
	v, _ := packet.GetBool()
	fmt.Println("true xor false =", v, ", passed data =", *data.(*int))
	soletta.Quit()
}

func Example_singleFlow() {
	soletta.Init()

	t := soletta.NewFlowNodeType("boolean/xor")
	pi1, _ := t.GetPort("IN[0]", soletta.FlowPortInput)
	pi2, _ := t.GetPort("IN[1]", soletta.FlowPortInput)
	po1, _ := t.GetPort("OUT", soletta.FlowPortOutput)
	inputPorts := []uint16{pi1, pi2}
	outputPorts := []uint16{po1}

	data := 11
	n := soletta.NewSingleFlowNode("boolean", *t, inputPorts, outputPorts, nil, singleFlowProcessCallback, &data)

	n.SendPacket("Bool", true, inputPorts[0])
	n.SendPacket("Bool", false, inputPorts[1])

	soletta.Run()

	n.Destroy()

	soletta.Shutdown()
	//Output: true xor false = true , passed data = 11
}
