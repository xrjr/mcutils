package ping

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"time"
	"unicode/utf16"

	"github.com/xrjr/mcutils/pkg/networking"
)

var (
	CommonLegacyRequest           [2]byte  = [2]byte{0xFE, 0x01}
	PluginMessagePacketIdentifier byte     = 0xFA
	SingleByteIdentifierValue     byte     = 0xFF
	Post1_3Padding                [6]byte  = [6]byte{0x00, 0xA7, 0x00, 0x31, 0x00, 0x00}
	Post1_3Delimiter              [2]byte  = [2]byte{0x00, 0x00}
	Pre1_3Delimiter               [2]byte  = [2]byte{0x00, 0xA7}
	MCPingHostStringWithLength    [24]byte = [24]byte{0x00, 0x0B, 0x00, 0x4D, 0x00, 0x43, 0x00, 0x7C, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x73, 0x00, 0x74} // MC|PingHost string, UTF16BE encoded, preceded by its length in characters as a short
	ProtocolNumber1_6_4           byte     = 78
)

// bigEndianUTF16ToString converts an array of byte containing an UTF-16BE encoded string into a string.
func bigEndianUTF16ToString(s []byte) string {
	if len(s)%2 != 0 {
		return ""
	}

	u16s := make([]uint16, 0, len(s)/2)

	for i := 0; i < len(s); i += 2 {
		u16s = append(u16s, binary.BigEndian.Uint16(s[i:i+2]))
	}

	return string(utf16.Decode(u16s))
}

// stringToBigEndianUTF16 converts a string into an array of byte containing an UTF-16BE encoded string.
func stringToBigEndianUTF16(s string) []byte {
	u16s := utf16.Encode([]rune(s))

	u8s := make([]byte, 0, len(u16s)*2)

	var buf [2]byte
	for _, u16 := range u16s {
		binary.BigEndian.PutUint16(buf[:], u16)
		u8s = append(u8s, buf[:]...)
	}

	return u8s
}

// generateLegacyPingRequest generates a networking.Output corresponding to a legacy ping request.
func generateLegacyPingRequest(hostname string, port uint16, use1_6_4protocol bool) networking.Output {
	out := networking.NewOutput()

	out.WriteBytes(CommonLegacyRequest[:])

	if use1_6_4protocol {
		out.WriteByte(PluginMessagePacketIdentifier)

		out.WriteBytes(MCPingHostStringWithLength[:])

		utf16hostname := stringToBigEndianUTF16(hostname)

		out.WriteBigEndianInt16(uint16(7 + len(utf16hostname)))

		out.WriteByte(ProtocolNumber1_6_4)

		out.WriteBigEndianInt16(uint16(len(hostname)))

		out.WriteBytes(utf16hostname)

		out.WriteBigEndianInt32(uint32(port))
	}

	return out
}

// parseLegacyPingResponse reads and parses a response (of type legacy ping) into a legacyPingResponse
func parseLegacyPingResponse(in networking.Input) (*legacyPingResponse, error) {
	var lpRes legacyPingResponse

	sbi, err := in.ReadByte()
	if err != nil {
		return nil, err
	}
	lpRes.SingleByteIdentifier = sbi

	if lpRes.SingleByteIdentifier != SingleByteIdentifierValue {
		return nil, ErrInvalidPacketType
	}

	length, err := in.ReadBigEndianInt16()
	if err != nil {
		return nil, err
	}
	lpRes.Length = length

	raw, err := in.ReadBytes(int(lpRes.Length) * 2)
	if err != nil {
		return nil, err
	}

	var delimiter []byte
	var post1_3 bool

	if len(raw) > 6 && bytes.Equal(raw[:6], Post1_3Padding[:]) {
		post1_3 = true
		delimiter = Post1_3Delimiter[:]
	} else {
		delimiter = Pre1_3Delimiter[:]
	}

	rawUTF16BEStrings := make([][]byte, 0, 5)
	current := 0
	rawUTF16BEStrings = append(rawUTF16BEStrings, []byte{})

	// If post 1.3, we start at 6th byte to skip Post1_3Padding
	padding := 0
	if post1_3 {
		padding = 6
	}

	for i := padding; i < int(lpRes.Length)*2; i += 2 {
		if bytes.Equal(raw[i:i+2], delimiter) {
			current++
			rawUTF16BEStrings = append(rawUTF16BEStrings, []byte{})
		} else {
			rawUTF16BEStrings[current] = append(rawUTF16BEStrings[current], raw[i:i+2]...)
		}
	}

	if post1_3 && len(rawUTF16BEStrings) != 5 {
		return nil, ErrMalformedPacket
	} else if !post1_3 && len(rawUTF16BEStrings) != 3 {
		return nil, ErrMalformedPacket
	}

	// common to pre and post 1.3
	lpRes.MOTD = bigEndianUTF16ToString(rawUTF16BEStrings[len(rawUTF16BEStrings)-3])

	online, err := strconv.Atoi(bigEndianUTF16ToString(rawUTF16BEStrings[len(rawUTF16BEStrings)-2]))
	if err != nil {
		return nil, err
	}
	lpRes.OnlinePlayers = online

	max, err := strconv.Atoi(bigEndianUTF16ToString(rawUTF16BEStrings[len(rawUTF16BEStrings)-1]))
	if err != nil {
		return nil, err
	}
	lpRes.MaxPlayers = max

	// post 1.3 specific
	if post1_3 {
		protocolVersion, err := strconv.Atoi(bigEndianUTF16ToString(rawUTF16BEStrings[len(rawUTF16BEStrings)-5]))
		if err != nil {
			return nil, err
		}
		lpRes.ProtocolVersion = protocolVersion

		lpRes.MinecraftVersion = bigEndianUTF16ToString(rawUTF16BEStrings[len(rawUTF16BEStrings)-4])
	}

	return &lpRes, nil
}

// PingClientLegacy is the legacy ping client.
type PingClientLegacy struct {
	hostname    string
	port        int
	dialOptions networking.DialTCPOptions
	conn        *networking.TCPConn
}

// NewClientLegacy returns a well-formed *LegacyPingClient.
func NewClientLegacy(hostname string, port int) *PingClientLegacy {
	return &PingClientLegacy{
		hostname: hostname,
		port:     port,
	}
}

// SetDialOptions sets the options used in the dial process of the connection.
func (client *PingClientLegacy) SetDialOptions(dialOptions networking.DialTCPOptions) {
	client.dialOptions = dialOptions
}

// Connect establishes a connection via TCP.
func (client *PingClientLegacy) Connect() error {
	if client.conn != nil {
		return networking.ErrConnectionAlreadyEstablished
	}

	conn, err := networking.DialTCP(client.hostname, client.port, client.dialOptions)
	if err != nil {
		return err
	}

	client.conn = conn
	return nil
}

// Ping sends a legacy ping request to the server, and returns various informations about the server, and the latency in ms.
// If the minecraft server has a version <= 1.3, ProtocolNumber and MinecraftVersion are not set.
func (client *PingClientLegacy) Ping() (LegacyPingInfos, int, error) {
	return client.ping(false)
}

// Ping1_6_4 sends a legacy ping request to the server (using 1.6+ SLP protocol), and returns various informations about the server, and the latency in ms.
// If the minecraft server has a version <= 1.3, ProtocolNumber and MinecraftVersion are not set.
func (client *PingClientLegacy) Ping1_6_4() (LegacyPingInfos, int, error) {
	return client.ping(true)
}

// ping sends a legacy ping request to the server, and returns various informations about the server, and the latency in ms.
func (client *PingClientLegacy) ping(use1_6_4protocol bool) (LegacyPingInfos, int, error) {
	if client.conn == nil {
		return LegacyPingInfos{}, -1, networking.ErrConnectionNotEstablished
	}

	pingRequest := generateLegacyPingRequest(client.hostname, uint16(client.port), use1_6_4protocol)

	start := time.Now().UnixMilli()
	pingResponse, err := client.conn.Send(pingRequest)
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}

	lpr, err := parseLegacyPingResponse(pingResponse)
	if err != nil {
		return LegacyPingInfos{}, -1, err
	}
	latency := time.Now().UnixMilli() - start

	return lpr.legacyPingInfos(), int(latency), nil
}

// Disconnect closes the connection.
// Connection is made not usable anymore no matter if the it closed properly or not.
func (client *PingClientLegacy) Disconnect() error {
	if client.conn == nil {
		return networking.ErrConnectionNotEstablished
	}

	err := client.conn.Close()
	client.conn = nil
	return err
}
