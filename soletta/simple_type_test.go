package soletta

import "strconv"

func simpleTypeProcessEvent(node *FlowNode, event *SimpleFlowEvent) bool {
	switch event.Type {
	case SimpleEventOpen:
		node.SetData(map[string]int32{"hour": -1, "minute": -1, "second": -1})
	case SimpleEventProcessInputPort:
		/* Output current time (gathered from input ports) as a string to output port */
		m := node.GetData().(map[string]int32)
		m[event.PortName], _ = event.Packet.GetInteger()

		if event.PortName == "second" && m["second"] != 0 && m["hour"] != -1 && m["minute"] != -1 {
			po, _ := node.GetPort("time", FlowPortOutput)
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
	Init()

	ports := []PortDescription{
		PortDescription{"hour", "Integer", FlowPortInput},
		PortDescription{"minute", "Integer", FlowPortInput},
		PortDescription{"second", "Integer", FlowPortInput},
		PortDescription{"time", "String", FlowPortOutput},
	}
	st := NewSimpleNodeType("custom", ports, simpleTypeProcessEvent)

	/* Create a flow network around the custom type
	   Redirect hour, minute and second data from wallclock type into
	   custom type input ports. Redirect custom type output port to console */
	b := NewFlowBuilder()
	b.AddNodeByTypeName("hour", "wallclock/hour", map[string]string{"send_initial_packet": "true"})
	b.AddNodeByTypeName("minute", "wallclock/minute", map[string]string{"send_initial_packet": "true"})
	b.AddNodeByTypeName("second", "wallclock/second", map[string]string{"send_initial_packet": "true"})
	b.AddNode("custom", st, map[string]string{"prefix": "", "suffix": ""})
	b.AddNodeByTypeName("console", "console", nil)
	b.Connect("hour", "OUT", "custom", "hour")
	b.Connect("minute", "OUT", "custom", "minute")
	b.Connect("second", "OUT", "custom", "second")
	b.Connect("custom", "time", "console", "IN")

	t := b.GetNodeType()
	defer t.Destroy()

	flow := t.CreateNode(nil, "highlevel", nil)

	Run()

	flow.Destroy()

	Shutdown()
}
