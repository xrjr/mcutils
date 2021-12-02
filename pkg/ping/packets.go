package ping

import "github.com/xrjr/mcutils/pkg/networking"

// packet is the common structure conatained in all ping packets
type packet struct {
	Length   uint32
	PacketID uint32
}

// handshakeResponse is the type respresenting the response of the handshake request
type handshakeResponse struct {
	packet
	JSONResponse map[string]interface{}
}

func (hsr *handshakeResponse) handshake() Handshake {
	return Handshake{
		Properties: hsr.JSONResponse,
	}
}

type Handshake struct {
	Properties map[string]interface{} `json:"properties"`
}

type pongResponse struct {
	packet
	Payload int64
}

// transformToPacket transforms an output payload into a readable ping packet, adding the length at the start and the packet id
func transformToPacket(out networking.Output) networking.Output {
	var packetOut networking.Output = networking.NewOutput()

	packetOut.WriteUVarInt(uint64(len(out.Bytes())))
	packetOut.WriteBytes(out.Bytes())

	return packetOut
}

// emptyPacket generates an empty packet ready to be sent
func emptyPacket(packetID uint32) networking.Output {
	out := networking.NewOutput()
	out.WriteUVarInt(uint64(packetID))
	return transformToPacket(out)
}
