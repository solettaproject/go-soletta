package soletta_test

import "testing"
import "github.com/solettaproject/go-soletta/soletta"
import "strconv"

func simpleTypeProcessEvent(node *soletta.FlowNode, event *soletta.SimpleFlowEvent) bool {
	switch event.Type {
	case soletta.SimpleEventOpen:
		node.SetData(map[string]int32{"hour": -1, "minute": -1, "second": -1})
	case soletta.SimpleEventProcessInputPort:
		/* Output current time (gathered from input ports) as a string to output port */
		m := node.GetData().(map[string]int32)
		m[event.PortName], _ = event.Packet.GetInteger()

		if event.PortName == "second" && m["second"] != 0 && m["hour"] != -1 && m["minute"] != -1 {
			po, _ := node.GetPort("time", soletta.FlowPortOutput)
			s := "Current time: " + strconv.Itoa(int(m["hour"])) + ":" + strconv.Itoa(int(m["minute"])) + ":" + strconv.Itoa(int(m["second"]))
			node.SendPacket("String", s, po)
		}
	}
	return true
}

/*
Shows the creation and usage of a custom type with
3 input ports of type Integer and 1 output port of type String
*/
func Example_simpleType() {
	soletta.Init()

	ports := []soletta.PortDescription{
		soletta.PortDescription{"hour", "Integer", soletta.FlowPortInput},
		soletta.PortDescription{"minute", "Integer", soletta.FlowPortInput},
		soletta.PortDescription{"second", "Integer", soletta.FlowPortInput},
		soletta.PortDescription{"time", "String", soletta.FlowPortOutput},
	}
	st := soletta.NewSimpleNodeType("custom", ports, simpleTypeProcessEvent)

	/* Create a flow network around the custom type
	   Redirect hour, minute and second data from wallclock type into
	   custom type input ports. Redirect custom type output port to console */
	b := soletta.NewFlowBuilder()
	b.AddNodeByTypeName("hour", "wallclock/hour", map[string]string{"send_initial_packet": "true"})
	b.AddNodeByTypeName("minute", "wallclock/minute", map[string]string{"send_initial_packet": "true"})
	b.AddNodeByTypeName("second", "wallclock/second", map[string]string{"send_initial_packet": "true"})
	b.AddNode("custom", st, map[string]string{"prefix": "", "suffix": ""})
	b.AddNodeByTypeName("console", "console", nil)
	b.Connect("hour", "OUT", -1, "custom", "hour", -1)
	b.Connect("minute", "OUT", -1, "custom", "minute", -1)
	b.Connect("second", "OUT", -1, "custom", "second", -1)
	b.Connect("custom", "time", -1, "console", "IN", -1)

	t := b.GetNodeType()
	defer t.Destroy()

	flow := t.CreateNode(nil, "highlevel", nil)

	soletta.Run()

	flow.Destroy()

	soletta.Shutdown()
}

func TestSimpleType(test *testing.T) {
	result := ""

	soletta.Init()

	processEvent := func(node *soletta.FlowNode, event *soletta.SimpleFlowEvent) bool {
		switch event.Type {
		case soletta.SimpleEventOpen:
			node.SetData(event.Options["prefix"])
		case soletta.SimpleEventProcessInputPort:
			s, _ := event.Packet.GetString()
			node.SendPacket("String", node.GetData().(string)+s, 0)
		}
		return true
	}
	ports := []soletta.PortDescription{
		soletta.PortDescription{"IN", "String", soletta.FlowPortInput},
		soletta.PortDescription{"OUT", "String", soletta.FlowPortOutput},
	}
	st := soletta.NewSimpleNodeType("custom", ports, processEvent)

	singleFlowProcessCallback := func(node *soletta.FlowNode, port uint16, packet *soletta.FlowPacket, data interface{}) {
		result, _ = packet.GetString()
		soletta.Quit()
	}
	n := soletta.NewSingleFlowNode("singleFlow", *st, []uint16{0}, []uint16{0}, map[string]string{"prefix": "Prefix: "}, singleFlowProcessCallback, nil)

	n.SendPacket("String", "Test string", 0)

	soletta.Run()

	n.Destroy()

	soletta.Shutdown()

	if result != "Prefix: Test string" {
		test.Fail()
	}
}
