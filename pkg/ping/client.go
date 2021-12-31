package ping

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/xrjr/mcutils/pkg/networking"
)

const (
	UnknownProtocolVersion int32  = -1
	HandshakePacketID      uint32 = 0
	PingPacketID           uint32 = 1
)

var (
	ErrInvalidPacketType error = errors.New("invalid packet type")
)

// generateHandshakeRequest generates a networking.Output corresponding to a handshake request.
func generateHandshakeRequest(hostname string, port uint16) networking.Output {
	out := networking.NewOutput()

	out.WriteUVarInt(uint64(HandshakePacketID))

	out.WriteVarInt(int64(UnknownProtocolVersion))

	out.WriteString(hostname)

	out.WriteBigEndianInt16(port)

	out.WriteUVarInt(1)

	return out
}

// parseHandshakeResponse reads and parses a response (of type handshake) into a handshakeResponse.
func parseHandshakeResponse(in networking.Input) (*handshakeResponse, error) {
	var hsRes handshakeResponse

	length, err := in.ReadUVarInt()
	if err != nil {
		return nil, err
	}
	hsRes.Length = uint32(length)

	packetID, err := in.ReadUVarInt()
	if err != nil {
		return nil, err
	}
	hsRes.PacketID = uint32(packetID)
	if hsRes.PacketID != HandshakePacketID {
		return nil, ErrInvalidPacketType
	}

	rawJSONResponse, err := in.ReadString()
	if err != nil {
		return nil, err
	}

	jsonResponse := make(map[string]interface{})
	err = json.Unmarshal([]byte(rawJSONResponse), &jsonResponse)
	if err != nil {
		return nil, err
	}

	hsRes.JSONResponse = jsonResponse

	return &hsRes, nil
}

// generatePingRequest generates a networking.Output corresponding to a ping request.
func generatePingRequest() networking.Output {
	out := networking.NewOutput()

	out.WriteUVarInt(uint64(PingPacketID))

	out.WriteBigEndianInt64(uint64(time.Now().UnixMilli()))

	return out
}

// parsePongResponse reads and parses a response (of type pong) into a *pongResponse.
func parsePongResponse(in networking.Input) (*pongResponse, error) {
	var pongRes pongResponse

	length, err := in.ReadUVarInt()
	if err != nil {
		return nil, err
	}
	pongRes.Length = uint32(length)

	packetID, err := in.ReadUVarInt()
	if err != nil {
		return nil, err
	}
	pongRes.PacketID = uint32(packetID)
	if pongRes.PacketID != PingPacketID {
		return nil, ErrInvalidPacketType
	}

	payload, err := in.ReadBigEndianInt64()
	if err != nil {
		return nil, err
	}
	pongRes.Payload = int64(payload)

	return &pongRes, nil
}

// PingClient is the ping client.
type PingClient struct {
	hostname string
	port     int
	conn     *networking.TCPConn

	// options
	SkipSRVLookup bool
	DialTimeout   time.Duration // in ms
	ReadTimeout   time.Duration // in ms
}

// NewClient returns a well-formed *PingClient.
func NewClient(hostname string, port int) *PingClient {
	return &PingClient{
		hostname: hostname,
		port:     port,

		SkipSRVLookup: false,
		DialTimeout:   5 * time.Second,
		ReadTimeout:   5 * time.Second,
	}
}

// Connect establishes a connection via TCP.
func (client *PingClient) Connect() error {
	if client.conn != nil {
		return networking.ErrConnectionAlreadyEstablished
	}

	conn, err := networking.DialTCP(client.hostname, client.port, networking.DialTCPOptions{
		SkipSRVLookup: client.SkipSRVLookup,
		DialTimeout:   client.DialTimeout,
	})
	if err != nil {
		return err
	}

	client.conn = conn
	return nil
}

// Handshake sends a handshake request to the server, and returns the formatted result.
func (client *PingClient) Handshake() (Handshake, error) {
	if client.conn == nil {
		return Handshake{}, networking.ErrConnectionNotEstablished
	}

	hsRequest := generateHandshakeRequest(client.hostname, uint16(client.port))
	hsRequestPacket := transformToPacket(hsRequest)
	fullHsRequest := networking.MergeOutputs(hsRequestPacket, emptyPacket(0))

	hsResponse, err := client.conn.Send(fullHsRequest)
	if err != nil {
		return Handshake{}, err
	}

	client.conn.SetReadDeadline(client.ReadTimeout)
	hs, err := parseHandshakeResponse(hsResponse)
	if err != nil {
		return Handshake{}, err
	}

	return hs.handshake(), nil
}

// Ping sends a ping request to the server, and returns the latency in ms.
// A ping request must be done after a handshake request has already been done.
func (client *PingClient) Ping() (int, error) {
	if client.conn == nil {
		return -1, networking.ErrConnectionNotEstablished
	}

	pingRequest := generatePingRequest()
	pingRequestPacket := transformToPacket(pingRequest)

	pingResponse, err := client.conn.Send(pingRequestPacket)
	if err != nil {
		return -1, err
	}

	client.conn.SetReadDeadline(client.ReadTimeout)
	pong, err := parsePongResponse(pingResponse)
	if err != nil {
		return -1, err
	}

	return int(time.Now().UnixMilli() - pong.Payload), nil
}

// Disconnect closes the connection.
// Connection is made not usable anymore no matter if the it closed properly or not.
func (client *PingClient) Disconnect() error {
	if client.conn == nil {
		return networking.ErrConnectionNotEstablished
	}

	err := client.conn.Close()
	client.conn = nil
	return err
}
