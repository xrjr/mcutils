package ping

import (
	"errors"

	"github.com/xrjr/mcutils/pkg/networking"
)

var (
	ErrMalformedPacket error = errors.New("malformed packet")
)

// JSON is an alias for map[string]interface{} that allows extracting data to an Infos struct.
type JSON map[string]interface{}

// packet is the common structure conatained in all ping packets.
type packet struct {
	Length   uint32
	PacketID uint32
}

// handshakeResponse is the type respresenting the response of the handshake request.
type handshakeResponse struct {
	packet
	JSONResponse map[string]interface{}
}

// handshake transforms the handshakeResponse into a more human-usable Handshake struct.
func (hsr *handshakeResponse) handshake() Handshake {
	return Handshake{
		Properties: hsr.JSONResponse,
	}
}

// Handshake contains SLP handshake informations.
type Handshake struct {
	Properties JSON `json:"properties"`
}

// pongResponse is the type respresenting the response of the ping request.
type pongResponse struct {
	packet
	Payload int64
}

// transformToPacket transforms an output payload into a readable ping packet, adding the length at the start and the packet id.
func transformToPacket(out networking.Output) networking.Output {
	var packetOut networking.Output = networking.NewOutput()

	packetOut.WriteUVarInt(uint64(len(out.Bytes())))
	packetOut.WriteBytes(out.Bytes())

	return packetOut
}

// emptyPacket generates an empty packet ready to be sent.
func emptyPacket(packetID uint32) networking.Output {
	out := networking.NewOutput()
	out.WriteUVarInt(uint64(packetID))
	return transformToPacket(out)
}

// legacyPingResponse is the type respresenting the response of the legacy ping request.
type legacyPingResponse struct {
	SingleByteIdentifier byte
	Length               uint16
	ProtocolVersion      int
	MinecraftVersion     string
	MOTD                 string
	OnlinePlayers        int
	MaxPlayers           int
}

// legacyPingInfos transforms the legacyPingResponse into a more human-usable LegacyPingInfos struct.
func (lpr *legacyPingResponse) legacyPingInfos() LegacyPingInfos {
	return LegacyPingInfos{
		ProtocolVersion:  lpr.ProtocolVersion,
		MinecraftVersion: lpr.MinecraftVersion,
		MOTD:             lpr.MOTD,
		OnlinePlayers:    lpr.OnlinePlayers,
		MaxPlayers:       lpr.MaxPlayers,
	}
}

// LegacyPingInfos contains legacy SLP informations.
type LegacyPingInfos struct {
	ProtocolVersion  int    `json:"protocolVersion"`
	MinecraftVersion string `json:"minecraftVersion"`
	MOTD             string `json:"motd"`
	OnlinePlayers    int    `json:"onlinePlayers"`
	MaxPlayers       int    `json:"maxPlayers"`
}
