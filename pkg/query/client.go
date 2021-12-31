package query

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/xrjr/mcutils/pkg/networking"
)

const (
	MagicValue uint16 = 65277
)

var (
	FullStatRequestPadding   = [4]byte{0x00, 0x00, 0x00, 0x00}
	FullStatResponsePadding1 = [11]byte{0x73, 0x70, 0x6C, 0x69, 0x74, 0x6E, 0x75, 0x6D, 0x00, 0x80, 0x00}
	FullStatResponsePadding2 = [10]byte{0x01, 0x70, 0x6C, 0x61, 0x79, 0x65, 0x72, 0x5F, 0x00, 0x00}
)

// generateSessionID generates a non-croptographically secure, non-seeded random valid session id.
func generateSessionID() uint32 {
	var res int32 = rand.Int31()
	return uint32(res) & 0x0F0F0F0F
}

// generateHandshakeRequest generates a networking.Output corresponding to a handshake request.
func generateHandshakeRequest(sessionID uint32) networking.Output {
	out := networking.NewOutput()

	out.WriteBigEndianInt16(MagicValue)

	out.WriteSingleByte(9)

	out.WriteBigEndianInt32(uint32(sessionID))

	return out
}

// parseHandshakeResponse reads and parses a response (of type handshake) into a handshakeResponse.
func parseHandshakeResponse(in networking.Input) (*handshakeResponse, error) {
	var hsRes *handshakeResponse = &handshakeResponse{}

	type_, err := in.ReadByte()
	if err != nil {
		return nil, err
	}
	hsRes.Type = type_

	sessionid, err := in.ReadBigEndianInt32()
	if err != nil {
		return nil, err
	}
	hsRes.SessionID = sessionid

	challengeTokenString, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	tmpIntChallengeToken, err := strconv.ParseInt(challengeTokenString, 10, 32)
	if err != nil {
		return nil, err
	}
	hsRes.ChallengeToken = uint32(tmpIntChallengeToken)

	return hsRes, nil
}

// generateBasicStatRequest generates a networking.Output corresponding to a basic stat query request.
func generateBasicStatRequest(sessionID, tokenID uint32) networking.Output {
	out := networking.NewOutput()

	out.WriteBigEndianInt16(MagicValue)

	out.WriteSingleByte(0)

	out.WriteBigEndianInt32(uint32(sessionID))

	out.WriteBigEndianInt32(uint32(tokenID))

	return out
}

// parseBasicStatResponse reads and parses a response (of type basic stat) into a *basicStatResponse.
func parseBasicStatResponse(in networking.Input) (*basicStatResponse, error) {
	var bsRes *basicStatResponse = &basicStatResponse{}

	type_, err := in.ReadByte()
	if err != nil {
		return nil, err
	}
	bsRes.Type = type_

	sessionid, err := in.ReadBigEndianInt32()
	if err != nil {
		return nil, err
	}
	bsRes.SessionID = sessionid

	motd, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	bsRes.MOTD = motd

	gametype, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	bsRes.GameType = gametype

	map_, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	bsRes.Map = map_

	numPlayersString, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	tmpIntNumPlayers, err := strconv.ParseInt(numPlayersString, 10, 64)
	if err != nil {
		return nil, err
	}
	bsRes.NumPlayers = int(tmpIntNumPlayers)

	maxPlayersString, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	tmpIntMaxPlayers, err := strconv.ParseInt(maxPlayersString, 10, 64)
	if err != nil {
		return nil, err
	}
	bsRes.MaxPlayers = int(tmpIntMaxPlayers)

	hostport, err := in.ReadLittleEndianInt16()
	if err != nil {
		return nil, err
	}
	bsRes.HostPort = int16(hostport)

	hostip, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	bsRes.HostIP = hostip

	return bsRes, nil
}

// generateFullStatRequest generates a networking.Request corresponding to a full stat query request.
func generateFullStatRequest(sessionID, tokenID uint32) networking.Output {
	out := networking.NewOutput()

	out.WriteBigEndianInt16(MagicValue)

	out.WriteSingleByte(0)

	out.WriteBigEndianInt32(uint32(sessionID))

	out.WriteBigEndianInt32(uint32(tokenID))

	out.WriteBytes(FullStatRequestPadding[:])

	return out
}

// parseFullStatResponse reads and parses a response (of type full stat) into a *fullStatResponse.
func parseFullStatResponse(in networking.Input) (*fullStatResponse, error) {
	var fsRes *fullStatResponse = &fullStatResponse{}

	type_, err := in.ReadByte()
	if err != nil {
		return nil, err
	}
	fsRes.Type = type_

	sessionid, err := in.ReadBigEndianInt32()
	if err != nil {
		return nil, err
	}
	fsRes.SessionID = sessionid

	padding1, err := in.ReadBytes(11)
	if err != nil {
		return nil, err
	}
	copy(fsRes.Padding1[:], padding1)

	fsRes.KVSection = make(map[string]string)
	var key string
	var value string
	for {
		key, err = in.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}

		if len(key) == 0 {
			break
		}

		value, err = in.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}

		fsRes.KVSection[key] = value
	}

	padding2, err := in.ReadBytes(10)
	if err != nil {
		return nil, err
	}
	copy(fsRes.Padding2[:], padding2)

	var player string
	for {
		player, err = in.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
		if len(player) == 0 {
			break
		}
		fsRes.PlayersSection = append(fsRes.PlayersSection, player)
	}

	return fsRes, nil
}

// QueryClient is the query client.
// challengeToken isn't stored in the client as it can change in the lifetime of a client, while session ID doesn't.
type QueryClient struct {
	hostname  string
	port      int
	conn      *networking.UDPConn
	sessionID uint32

	// options
	SkipSRVLookup                bool
	ForceUDPProtocolForSRVLookup bool
	DialTimeout                  time.Duration
	ReadTimeout                  time.Duration
}

// NewClient returns a well-formed *QueryClient.
func NewClient(hostname string, port int) *QueryClient {
	return &QueryClient{
		hostname: hostname,
		port:     port,

		SkipSRVLookup:                false,
		ForceUDPProtocolForSRVLookup: false,
		DialTimeout:                  5 * time.Second,
		ReadTimeout:                  5 * time.Second,
	}
}

// Connect establishes a connection via UDP.
func (client *QueryClient) Connect() error {
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

	client.sessionID = generateSessionID()
	client.conn = conn
	return nil
}

// Handshake sends a handshake query to the server, and returns the challenge token if successful.
func (client *QueryClient) Handshake() (uint32, error) {
	if client.conn == nil {
		return 0, networking.ErrConnectionNotEstablished
	}

	hsRequest := generateHandshakeRequest(client.sessionID)

	// UDPConn reads are made in send method
	err := client.conn.SetReadDeadline(client.ReadTimeout)
	if err != nil {
		return 0, err
	}

	hsResponse, err := client.conn.Send(hsRequest)
	if err != nil {
		return 0, err
	}

	hs, err := parseHandshakeResponse(hsResponse)
	if err != nil {
		return 0, err
	}

	return hs.ChallengeToken, nil
}

// BasicStat sends a basic stat query to the server, and returns the formatted result.
func (client *QueryClient) BasicStat(challengeToken uint32) (BasicStat, error) {
	if client.conn == nil {
		return BasicStat{}, networking.ErrConnectionNotEstablished
	}

	bsRequest := generateBasicStatRequest(client.sessionID, challengeToken)

	// UDPConn reads are made in send method
	err := client.conn.SetReadDeadline(client.ReadTimeout)
	if err != nil {
		return BasicStat{}, err
	}

	bsResponse, err := client.conn.Send(bsRequest)
	if err != nil {
		return BasicStat{}, err
	}

	bs, err := parseBasicStatResponse(bsResponse)
	if err != nil {
		return BasicStat{}, err
	}

	return bs.basicStat(), nil
}

// FullStat sends a full stat query to the server, and returns the formatted result.
func (client *QueryClient) FullStat(challengeToken uint32) (FullStat, error) {
	if client.conn == nil {
		return FullStat{}, networking.ErrConnectionNotEstablished
	}

	fsRequest := generateFullStatRequest(client.sessionID, challengeToken)

	// UDPConn reads are made in send method
	err := client.conn.SetReadDeadline(client.ReadTimeout)
	if err != nil {
		return FullStat{}, err
	}

	fsResponse, err := client.conn.Send(fsRequest)
	if err != nil {
		return FullStat{}, err
	}

	fs, err := parseFullStatResponse(fsResponse)
	if err != nil {
		return FullStat{}, err
	}

	return fs.fullStat(), nil
}

// Disconnect closes the connection.
// Connection is made not usable anymore no matter if the it closed properly or not.
func (client *QueryClient) Disconnect() error {
	if client.conn == nil {
		return networking.ErrConnectionNotEstablished
	}

	err := client.conn.Close()
	client.conn = nil
	return err
}
