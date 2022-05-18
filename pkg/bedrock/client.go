package bedrock

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/xrjr/mcutils/pkg/networking"
)

const (
	UnconnectedPing = 0x01
)

var (
	// MagicBedrockValue
	// hardcoded magic  https://github.com/facebookarchive/RakNet/blob/1a169895a900c9fc4841c556e16514182b75faf8/Source/RakPeer.cpp#L135
	MagicBedrockValue = []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78}

	ErrEmptyResponse        = errors.New("empty response")
	ErrNotIdUnconnectedPing = errors.New("first byte is not id_unconnected_pong")
	ErrNotMagicBytes        = errors.New("magic bytes do not match")
)

// BedrockClient raklib-query client for bedrock edition
type BedrockClient struct {
	hostname  string
	port      int
	conn      *networking.UDPConn
	sessionID []byte

	// options
	SkipSRVLookup                bool
	ForceUDPProtocolForSRVLookup bool
	DialTimeout                  time.Duration
	ReadTimeout                  time.Duration
}

// generateSessionID generates the session id
func generateSessionID() []byte {
	sid := make([]byte, 8)
	binary.LittleEndian.PutUint64(sid, 2)

	return sid
}

// generateTimestamp 64bit current time as bytes
func generateTimestamp() []byte {
	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(time.Now().Unix()))

	return t
}

// generateStatRequest generates a networking.Request corresponding un
func generateStatRequest(sessionID []byte) networking.Output {
	out := networking.NewOutput()

	out.WriteSingleByte(UnconnectedPing)

	out.WriteBytes(generateTimestamp())

	out.WriteBytes(MagicBedrockValue)

	out.WriteBytes(sessionID)

	return out
}

// parseStatResponse reads and parses a response into a *BEStat
func parseStatResponse(in networking.Input) (*BEStat, error) {
	response, err := in.ReadString()
	if err != nil {
		return nil, err
	}

	// TODO: if server-name or motd contains a ';' it is no escaped, and will break this parsing
	data := strings.Split(response, ";")

	if len(data) == 0 {
		return nil, ErrEmptyResponse
	}

	stat := &BEStat{}

	// Yes, it looks disgusting, but I think this is the best option to parse the response
	// since the response comes in a string with ';' as a divide

	if len(data) >= 1 {
		stat.GameName = data[0]
	}

	if len(data) >= 2 {
		stat.MOTD = data[1]
	}

	if len(data) >= 3 {
		stat.Protocol = data[2]
	}

	if len(data) >= 4 {
		stat.Version = data[3]
	}

	if len(data) >= 5 {
		stat.Players, _ = strconv.Atoi(data[4])
	}

	if len(data) >= 6 {
		stat.MaxPlayers, _ = strconv.Atoi(data[5])
	}

	if len(data) >= 7 {
		stat.ServerID, _ = strconv.ParseInt(data[6], 10, 64)
	}

	if len(data) >= 8 {
		stat.Map = data[7]
	}

	if len(data) >= 9 {
		stat.GameMode = data[8]
	}

	if len(data) >= 10 {
		stat.NintendoLimited = data[9]
	}

	if len(data) >= 11 {
		stat.IPv4Port, _ = strconv.Atoi(data[10])
	}

	if len(data) >= 12 {
		stat.IPv6Port, _ = strconv.Atoi(data[11])
	}

	if len(data) >= 13 {
		// What is this?
		stat.Extra = data[12]
	}

	return stat, err
}

// NewClient returns a formed *BedrockClient
func NewClient(hostname string, port int) *BedrockClient {
	return &BedrockClient{
		hostname: hostname,
		port:     port,

		SkipSRVLookup:                false,
		ForceUDPProtocolForSRVLookup: false,
		DialTimeout:                  5 * time.Second,
		ReadTimeout:                  5 * time.Second,
	}
}

// Connect establishes a connection via UDP.
func (c *BedrockClient) Connect() error {
	if c.conn != nil {
		return networking.ErrConnectionAlreadyEstablished
	}

	conn, err := networking.DialUDP(c.hostname, c.port, networking.DialUDPOptions{
		SkipSRVLookup:                c.SkipSRVLookup,
		ForceUDPProtocolForSRVLookup: c.ForceUDPProtocolForSRVLookup,
		DialTimeout:                  c.DialTimeout,
	})
	if err != nil {
		return err
	}

	c.sessionID = generateSessionID()
	c.conn = conn
	return nil
}

// Stat sends a request to get information about the server and returns the result.
func (c *BedrockClient) Stat() (*BEStat, error) {
	statRequest := generateStatRequest(c.sessionID)

	err := c.conn.SetReadDeadline(c.ReadTimeout)
	if err != nil {
		return nil, err
	}

	statResponse, err := c.conn.Send(statRequest)
	if err != nil {
		return nil, err
	}

	packetId, err := statResponse.ReadByte()
	if err != nil {
		return nil, err
	}

	if packetId != 0x1C { // 0x1C = ID_UNCONNECTED_PONG
		return nil, ErrNotIdUnconnectedPing
	}

	// 0-16 - ???
	// 16-32 - magic
	// 33 - ???
	magic, _ := statResponse.ReadBytes(33)
	if !bytes.Equal(magic[16:32], MagicBedrockValue) {
		return nil, ErrNotMagicBytes
	}

	return parseStatResponse(statResponse)
}

// Disconnect closes the connection.
// Connection is made not usable anymore no matter if the it closed properly or not.
func (c *BedrockClient) Disconnect() error {
	if c.conn == nil {
		return networking.ErrConnectionNotEstablished
	}

	err := c.conn.Close()
	c.conn = nil
	return err
}
