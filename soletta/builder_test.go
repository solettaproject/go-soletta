package soletta

func Example() {
	Init()

	b := NewFlowBuilder()
	b.AddNode("keyboard", "keyboard/int", nil)
	b.AddNode("console", "console", map[string]string{"prefix": "Hello: ", "suffix": " Bye"})
	b.Connect("keyboard", "OUT", "console", "IN")

	t := b.GetNodeType()
	defer t.Destroy()

	flow := t.CreateNode(nil, "highlevel", FlowOptions{})

	Run()

	flow.Destroy()

	Shutdown()
}
