package soletta_test

import "github.com/solettaproject/go-soletta/soletta"
import "testing"
import "strconv"

func Example_flowBuilder() {
	soletta.Init()

	b := soletta.NewFlowBuilder()
	b.AddNodeByTypeName("keyboard", "keyboard/int", nil)
	b.AddNodeByTypeName("console", "console", map[string]string{"prefix": "Hello: ", "suffix": " Bye"})
	b.Connect("keyboard", "OUT", -1, "console", "IN", -1)

	t := b.GetNodeType()
	defer t.Destroy()

	flow := t.CreateNode(nil, "highlevel", nil)

	soletta.Run()

	flow.Destroy()

	soletta.Shutdown()
}

func TestBuilder(test *testing.T) {
	result := true

	soletta.Init()

	/* Create a custom type with an output Integer port
	   that sends to it the value passed as option */
	customType1ProcessEvent := func(node *soletta.FlowNode, event *soletta.SimpleFlowEvent) bool {
		switch event.Type {
		case soletta.SimpleEventOpen:
			d, _ := strconv.Atoi(event.Options["value"])
			node.SendPacket("Integer", int32(d), 0)
		}
		return true
	}
	ports1 := []soletta.PortDescription{
		soletta.PortDescription{"OUT", "Integer", soletta.FlowPortOutput},
	}
	st1 := soletta.NewSimpleNodeType("custom1", ports1, customType1ProcessEvent)

	/* Create an equivalent of console node to catch output */
	customType2ProcessEvent := func(node *soletta.FlowNode, event *soletta.SimpleFlowEvent) bool {
		switch event.Type {
		case soletta.SimpleEventProcessInputPort:
			result, _ = event.Packet.GetBool()
			soletta.Quit()
		}
		return true
	}
	ports2 := []soletta.PortDescription{
		soletta.PortDescription{"IN", "Bool", soletta.FlowPortInput},
	}
	st2 := soletta.NewSimpleNodeType("custom2", ports2, customType2ProcessEvent)

	/* Create a flow network around int/equal and compare values 13 and 14 */
	b := soletta.NewFlowBuilder()
	b.AddNode("custom1Node1", st1, map[string]string{"value": "13"})
	b.AddNode("custom1Node2", st1, map[string]string{"value": "14"})
	b.AddNodeByTypeName("int", "int/equal", nil)
	b.AddNode("custom2", st2, nil)

	b.Connect("custom1Node1", "OUT", -1, "int", "IN", 0)
	b.Connect("custom1Node2", "OUT", -1, "int", "IN", 1)
	b.Connect("int", "OUT", -1, "custom2", "IN", -1)

	t := b.GetNodeType()
	defer t.Destroy()

	flow := t.CreateNode(nil, "highlevel", nil)

	soletta.Run()

	flow.Destroy()

	soletta.Shutdown()

	if result {
		test.Fail()
	}
}

func TestBuilder2(test *testing.T) {
	soletta.Init()

	b := soletta.NewFlowBuilder()
	b.AddNodeByTypeName("int", "int/addition", nil)
	b.ExportPort("int", "OPERAND", 0, "IN0", soletta.FlowPortInput)
	b.ExportPort("int", "OPERAND", 1, "IN1", soletta.FlowPortInput)
	b.ExportPort("int", "OPERAND", 2, "IN2", soletta.FlowPortInput)
	b.ExportPort("int", "OUT", -1, "OUT", soletta.FlowPortOutput)
	t := b.GetNodeType()
	defer t.Destroy()

	nip, nop := t.GetPortCount(soletta.FlowPortInput), t.GetPortCount(soletta.FlowPortOutput)
	if nip != 3 || nop != 1 {
		test.Fail()
	}
}
