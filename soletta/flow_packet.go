package soletta

/*
#include <soletta.h>
#include <sol-flow.h>
*/
import "C"
import "time"
import "image/color"
import "unsafe"

//Data structure describing a geographical location
type Location struct {
	Latitude, Longitude, Altitude float64
}

//Data structure implementing the image.Color interface
type Color struct {
	Red, Green, Blue, Alpha uint32
}

func (color Color) RGBA() (r, g, b, a uint32) {
	return color.Red, color.Green, color.Blue, color.Alpha
}

//Data structure describing a vector
type DirectionVector struct {
	X, Y, Z, Max, Min float64
}

//Represents data packets exchanged between flow nodes
type FlowPacket struct {
	cpacket *C.struct_sol_flow_packet
}

//Constructs and returns a new flow packet
func NewFlowPacket(name string, args ...interface{}) *FlowPacket {
	switch name {
	case "String":
		cs := C.CString(args[0].(string))
		defer C.free(unsafe.Pointer(cs))
		return &FlowPacket{C.sol_flow_packet_new_string(cs)}
	case "Integer":
		return &FlowPacket{C.sol_flow_packet_new_irange_value(C.int32_t(args[0].(int32)))}
	case "Double":
		return &FlowPacket{C.sol_flow_packet_new_drange_value(C.double(args[0].(float64)))}
	case "Bool":
		return &FlowPacket{C.sol_flow_packet_new_bool(C.bool(args[0].(bool)))}
	case "Byte":
		return &FlowPacket{C.sol_flow_packet_new_byte(C.uchar(args[0].(byte)))}
	case "Direction":
		dv := args[0].(DirectionVector)
		cdv := C.struct_sol_direction_vector{C.double(dv.Max), C.double(dv.Min), C.double(dv.X), C.double(dv.Y), C.double(dv.Z)}
		return &FlowPacket{C.sol_flow_packet_new_direction_vector(&cdv)}
	case "Color":
		r, g, b, _ := args[0].(color.Color).RGBA()
		return &FlowPacket{C.sol_flow_packet_new_rgb_components(C.uint32_t(r), C.uint32_t(g), C.uint32_t(b))}
	case "Location":
		loc := args[0].(Location)
		return &FlowPacket{C.sol_flow_packet_new_location_components(C.double(loc.Latitude), C.double(loc.Longitude), C.double(loc.Altitude))}
	case "Time":
		t := args[0].(time.Time)
		ctimespec := C.struct_timespec{C.__time_t(t.Unix()), C.__syscall_slong_t(t.UnixNano())}
		return &FlowPacket{C.sol_flow_packet_new_timestamp(&ctimespec)}
	}
	return nil
}

//Returns the integer value stored in the packet
func (fp *FlowPacket) GetInteger() (ret int32, ok bool) {
	var value C.int32_t
	r := C.sol_flow_packet_get_irange_value(fp.cpacket, &value)
	ret = int32(value)
	ok = r >= 0
	return ret, ok
}

//Returns the double value store in the packet
func (fp *FlowPacket) GetDouble() (ret float64, ok bool) {
	var value C.double
	r := C.sol_flow_packet_get_drange_value(fp.cpacket, &value)
	ret = float64(value)
	ok = r >= 0
	return ret, ok
}

//Returns the boolean value stored in the packet
func (fp *FlowPacket) GetBool() (ret, ok bool) {
	var value C.bool
	r := C.sol_flow_packet_get_bool(fp.cpacket, &value)
	ret = bool(value)
	ok = r >= 0
	return ret, ok
}

//Returns the byte value stored in the packet
func (fp *FlowPacket) GetByte() (ret byte, ok bool) {
	var value C.uchar
	r := C.sol_flow_packet_get_byte(fp.cpacket, &value)
	ret = byte(value)
	ok = r >= 0
	return ret, ok
}

//Returns the string value stored in the packet
func (fp *FlowPacket) GetString() (ret string, ok bool) {
	var value *C.char
	r := C.sol_flow_packet_get_string(fp.cpacket, &value)
	if r < 0 {
		return "", false
	}
	return C.GoString(value), true
}

//Returns the location value stored in the packet
func (fp *FlowPacket) GetLocation() (ret Location, ok bool) {
	var alt, lon, lat C.double
	r := C.sol_flow_packet_get_location_components(fp.cpacket, &lat, &lon, &alt)
	ret = Location{float64(alt), float64(lon), float64(alt)}
	ok = r >= 0
	return ret, ok
}

//Returns the time value stored in the packet
func (fp *FlowPacket) GetTime() (ret time.Time, ok bool) {
	var value C.struct_timespec
	r := C.sol_flow_packet_get_timestamp(fp.cpacket, &value)
	ret = time.Unix(int64(value.tv_sec), int64(value.tv_nsec))
	ok = r >= 0
	return ret, ok
}

//Returns the RGBA color stored in the packet
func (fp *FlowPacket) GetColor() (ret color.Color, ok bool) {
	var value C.struct_sol_rgb
	cr := C.sol_flow_packet_get_rgb(fp.cpacket, &value)
	r := uint32(float64(value.red) / float64(value.red_max) * 65536)
	g := uint32(float64(value.green) / float64(value.green_max) * 65536)
	b := uint32(float64(value.blue) / float64(value.blue_max) * 65536)
	ret = Color{r, g, b, 65536}
	ok = cr >= 0
	return ret, ok
}

//Returns the direction vector stored in the packet
func (fp *FlowPacket) GetDirection() (ret DirectionVector, ok bool) {
	var value C.struct_sol_direction_vector
	r := C.sol_flow_packet_get_direction_vector(fp.cpacket, &value)
	ret = DirectionVector{float64(value.x), float64(value.y), float64(value.z), float64(value.max), float64(value.min)}
	ok = r >= 0
	return ret, ok
}

//Frees the resources associated with the flow packet
func (fp *FlowPacket) Destroy() {
	C.sol_flow_packet_del(fp.cpacket)
}
