package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
*/
import "C"

//Represents a port used in flow node connections
type FlowPort struct {
}

const (
	FlowPortInput  int = iota
	FlowPortOutput int = iota
)

//Data structure that describes a port
type PortDescription struct {
	Name       string
	PacketType string
	PortType   int
}
