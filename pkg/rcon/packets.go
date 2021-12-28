package rcon

import "github.com/xrjr/mcutils/pkg/networking"

// packet is the structure representing an entire rcon packet, including the padding at the end.
type packet struct {
	Length    uint32
	RequestID int32
	Type      uint32
	Payload   string
	Padding   byte
}

// transformToPacket transforms an output payload into a readable rcon packet, adding the length at the start.
func transformToPacket(out networking.Output) networking.Output {
	var packetOut networking.Output = networking.NewOutput()

	packetOut.WriteLittleEndianInt32(uint32(len(out.Bytes())))
	packetOut.WriteBytes(out.Bytes())

	return packetOut
}
