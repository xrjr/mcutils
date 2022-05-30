package bedrock

import (
	"bytes"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/xrjr/mcutils/pkg/networking"
)

const (
	UnconnectedPingPacketID byte = 0x01
	UnconnectedPongPacketID byte = 0x1C
)

var (
	ErrInvalidPacketType error = errors.New("invalid packet type")
	ErrInvalidMagic      error = errors.New("invalid magic")
	ErrInvalidData       error = errors.New("invalid data")

	RaknetMagic = [16]byte{0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}
)

// generateUnconnectedPingRequest generates a networking.Output corresponding to an unconnected ping request.
func generateUnconnectedPingRequest(clientGUID uint64) networking.Output {
	out := networking.NewOutput()

	out.WriteByte(UnconnectedPingPacketID)

	out.WriteBigEndianInt64(uint64(time.Now().Unix()))

	out.WriteBytes(RaknetMagic[:])

	out.WriteBigEndianInt64(clientGUID)

	return out
}

// parseUnconnectedPongResponse reads and parses a response (of type unconncted pong) into a unconnectedPongResponse.
func parseUnconnectedPongResponse(in networking.Input) (*unconnectedPongResponse, error) {
	var res unconnectedPongResponse

	packetID, err := in.ReadByte()
	if err != nil {
		return nil, err
	}
	res.PacketID = packetID
	if res.PacketID != UnconnectedPongPacketID {
		return nil, ErrInvalidPacketType
	}

	clientTimestamp, err := in.ReadBigEndianInt64()
	if err != nil {
		return nil, err
	}
	res.ClientTimestamp = clientTimestamp

	serverGUID, err := in.ReadBigEndianInt64()
	if err != nil {
		return nil, err
	}
	res.ServerGUID = serverGUID

	magic, err := in.ReadBytes(16)
	if err != nil {
		return nil, err
	}
	res.Magic = magic
	if !bytes.Equal(res.Magic, RaknetMagic[:]) {
		return nil, ErrInvalidMagic
	}

	data, err := in.ReadRaknetString()
	if err != nil {
		return nil, err
	}

	splittedData := strings.Split(data, ";")

	if len(splittedData) >= 1 {
		res.GameName = splittedData[0]
	}

	if len(splittedData) >= 2 {
		res.MOTD = splittedData[1]

	}

	if len(splittedData) >= 3 && splittedData[2] != "" {
		res.ProtocolVersion, err = strconv.Atoi(splittedData[2])
		if err != nil {
			return nil, err
		}
	}

	if len(splittedData) >= 4 {
		res.MinecraftVersion = splittedData[3]

	}

	if len(splittedData) >= 5 && splittedData[4] != "" {
		res.OnlinePlayers, err = strconv.Atoi(splittedData[4])
		if err != nil {
			return nil, err
		}
	}

	if len(splittedData) >= 6 && splittedData[5] != "" {
		res.MaxPlayers, err = strconv.Atoi(splittedData[5])
		if err != nil {
			return nil, err
		}
	}

	if len(splittedData) >= 7 {
		res.ServerID = splittedData[6]

	}

	if len(splittedData) >= 8 {
		res.LevelName = splittedData[7]

	}

	if len(splittedData) >= 9 {
		res.GameMode = splittedData[8]

	}

	if len(splittedData) >= 10 && splittedData[9] != "" {
		res.GameModeNumeric, err = strconv.Atoi(splittedData[9])
		if err != nil {
			return nil, err
		}
	}

	if len(splittedData) >= 11 && splittedData[10] != "" {
		res.IPv4Port, err = strconv.Atoi(splittedData[10])
		if err != nil {
			return nil, err
		}
	}

	if len(splittedData) >= 12 && splittedData[11] != "" {
		res.IPv6Port, err = strconv.Atoi(splittedData[11])
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

// PingClient is the bedrock ping client.
type PingClient struct {
	hostname   string
	port       int
	conn       *networking.UDPConn
	ClientGUID uint64

	// options
	SkipSRVLookup                bool
	ForceUDPProtocolForSRVLookup bool
	DialTimeout                  time.Duration
	ReadTimeout                  time.Duration
}

// NewClient returns a well-formed *PingClient.
// ClientGUID is set to a random value.
func NewClient(hostname string, port int) *PingClient {
	return &PingClient{
		hostname:   hostname,
		port:       port,
		ClientGUID: uint64(rand.Int()),

		SkipSRVLookup:                false,
		ForceUDPProtocolForSRVLookup: false,
		DialTimeout:                  5 * time.Second,
		ReadTimeout:                  5 * time.Second,
	}
}

// Connect establishes a connection via UDP.
func (client *PingClient) Connect() error {
	if client.conn != nil {
		return networking.ErrConnectionAlreadyEstablished
	}

	conn, err := networking.DialUDP(client.hostname, client.port, networking.DialUDPOptions{
		SkipSRVLookup:                client.SkipSRVLookup,
		ForceUDPProtocolForSRVLookup: client.ForceUDPProtocolForSRVLookup,
		DialTimeout:                  client.DialTimeout,
	})
	if err != nil {
		return err
	}

	client.conn = conn
	return nil
}

// Handshake sends an unconncted ping request to the server, and returns the pong response informations.
func (client *PingClient) UnconnectedPing() (UnconnectedPong, int, error) {
	if client.conn == nil {
		return UnconnectedPong{}, -1, networking.ErrConnectionNotEstablished
	}

	unconnectedPingRequest := generateUnconnectedPingRequest(client.ClientGUID)

	startTime := time.Now().UnixMilli()

	// UDPConn reads are made in send method
	err := client.conn.SetReadDeadline(client.ReadTimeout)
	if err != nil {
		return UnconnectedPong{}, -1, err
	}

	unconnectedPongResponse, err := client.conn.Send(unconnectedPingRequest)
	if err != nil {
		return UnconnectedPong{}, -1, err
	}

	unconnectedPong, err := parseUnconnectedPongResponse(unconnectedPongResponse)
	if err != nil {
		return UnconnectedPong{}, -1, err
	}

	return unconnectedPong.unconnectedPong(), int(time.Now().UnixMilli() - startTime), nil
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
