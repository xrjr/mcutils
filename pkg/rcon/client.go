package rcon

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"

	"github.com/xrjr/mcutils/pkg/networking"
)

const (
	PacketSizeEmptyPayload       int = 14
	MaximumResponsePayloadLength int = 4096
	MaximumRequestPayloadLength  int = 1446
	LoginRequestType             int = 3
	CommandRequestType           int = 2
	WrongPasswordResponseType    int = 2
	CommandResponseType          int = 0
	InvalidRequestType           int = 4
)

var (
	ErrNotAuthenticated error = errors.New("authentication hasn't been realised or has failed. Call Authenticate method to authenticate")
	ErrCommandTooLong   error = errors.New("command length must be 1446 or less")
)

// generateRequestID generates a non-croptographically secure, non-seeded random request id.
func generateRequestID() uint32 {
	var res int32 = rand.Int31()
	return uint32(res)
}

// generateLoginRequest generates a networking.Output corresponding to a login request.
func generateLoginRequest(requestID uint32, password string) networking.Output {
	out := networking.NewOutput()

	out.WriteLittleEndianInt32(requestID)

	out.WriteLittleEndianInt32(uint32(LoginRequestType))

	out.WriteNullTerminatedString(password)

	out.WriteSingleByte(0)

	return out
}

// generateCommandRequest generates a networking.Output corresponding to a command request.
func generateCommandRequest(requestID uint32, command string) networking.Output {
	out := networking.NewOutput()

	out.WriteLittleEndianInt32(requestID)

	out.WriteLittleEndianInt32(uint32(CommandRequestType))

	out.WriteNullTerminatedString(command)

	out.WriteSingleByte(0)

	return out
}

// generateInvalidRequest generates an invalid request, that will not be understood by the server.
// It is useful to defragment multi-packet response.
func generateInvalidRequest(requestID uint32) networking.Output {
	out := networking.NewOutput()

	out.WriteLittleEndianInt32(requestID)

	out.WriteLittleEndianInt32(uint32(InvalidRequestType))

	out.WriteNullTerminatedString("")

	out.WriteSingleByte(0)

	return out
}

// parsePacket reads and parses an input into a *packet.
func parsePacket(in networking.Input) (*packet, error) {
	var p packet

	length, err := in.ReadLittleEndianInt32()
	if err != nil {
		return nil, err
	}
	p.Length = length

	requestID, err := in.ReadLittleEndianInt32()
	if err != nil {
		return nil, err
	}
	p.RequestID = int32(requestID)

	type_, err := in.ReadLittleEndianInt32()
	if err != nil {
		return nil, err
	}
	p.Type = type_

	payload, err := in.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}
	p.Payload = payload

	padding, err := in.ReadByte()
	if err != nil {
		return nil, err
	}
	p.Padding = padding

	return &p, nil
}

// parseFragmentedPacket reads and parse a response into a packet, and return its payload as a string.
// If the payload is fragmented into multiple packets, it will send and invalid packet to the server and read all of the packets received (including the invalid response), to defragment the multi-packet response.
// In this case, packet returned is the first recieved, and the payload return is the concatenation of all the payloads.
func parseFragmentedPacketPayload(conn networking.TCPConn, in networking.Input) (*packet, string, error) {
	var buf *bytes.Buffer = &bytes.Buffer{}
	p, err := parsePacket(in)
	if err != nil {
		return nil, "", err
	}
	buf.WriteString(p.Payload)

	if len(p.Payload) >= MaximumResponsePayloadLength {
		invalidRequest := generateInvalidRequest(uint32(p.RequestID))
		invalidRequestPacket := transformToPacket(invalidRequest)
		_, err = conn.Send(invalidRequestPacket) // response is the same as the one transmitted in param
		if err != nil {
			return nil, "", err
		}

		var tmpPacket *packet
		for {
			tmpPacket, err = parsePacket(in)
			if err != nil {
				return nil, "", err
			}

			if tmpPacket.Payload == fmt.Sprintf("Unknown request %x", InvalidRequestType) {
				break
			}

			if tmpPacket.RequestID == p.RequestID {
				buf.WriteString(tmpPacket.Payload)
			}
		}
	}

	return p, buf.String(), nil
}

// RCONClient is the RCON client.
type RCONClient struct {
	hostname      string
	port          int
	dialOptions   networking.DialTCPOptions
	conn          *networking.TCPConn
	authenticated bool
}

// NewClient returns a well-formed *RCONClient.
func NewClient(hostname string, port int) *RCONClient {
	return &RCONClient{
		hostname: hostname,
		port:     port,
	}
}

// SetDialOptions sets the options used in the dial process of the connection.
func (client *RCONClient) SetDialOptions(dialOptions networking.DialTCPOptions) {
	client.dialOptions = dialOptions
}

// Connect establishes a connection via TCP.
func (client *RCONClient) Connect() error {
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

// Authenticate sends an authentication packet to the server.
// If connection is successful it returns true, if not it returns false.
// If the communication didn't go wrong, err will be nil. This means that (false, nil) is a perfectly fine return if the password was wrong but he communication went well.
func (client *RCONClient) Authenticate(password string) (bool, error) {
	if client.conn == nil {
		return false, networking.ErrConnectionNotEstablished
	}

	rid := generateRequestID()

	loginRequest := generateLoginRequest(rid, password)

	loginRequestPacket := transformToPacket(loginRequest)

	loginResponse, err := client.conn.Send(loginRequestPacket)
	if err != nil {
		return false, err
	}

	packet, err := parsePacket(loginResponse)
	if err != nil {
		return false, err
	}

	if packet.RequestID == -1 {
		return false, nil
	}

	client.authenticated = true
	return true, nil
}

// Command sends a command packet to the server.
// Command length cannot be over MaximumRequestPayloadLength. This is a limitation of the source RCON protocol.
func (client *RCONClient) Command(command string) (string, error) {
	if len(command) > MaximumRequestPayloadLength {
		return "", ErrCommandTooLong
	}
	if client.conn == nil {
		return "", networking.ErrConnectionNotEstablished
	}
	if !client.authenticated {
		return "", ErrNotAuthenticated
	}

	rid := generateRequestID()

	commandRequest := generateCommandRequest(rid, command)

	commandRequestPacket := transformToPacket(commandRequest)

	commandResponse, err := client.conn.Send(commandRequestPacket)
	if err != nil {
		return "", err
	}

	packet, defragementedPayload, err := parseFragmentedPacketPayload(*client.conn, commandResponse)
	if err != nil {
		return "", err
	}

	if packet.RequestID == -1 {
		return "", ErrNotAuthenticated
	}

	return defragementedPayload, nil
}

// Disconnect closes the connection. This also means the authentication isn't valid anymore, even if a call to the Connect method is made.
// Connection is made not usable anymore no matter if the it closed properly or not.
func (client *RCONClient) Disconnect() error {
	client.authenticated = false

	if client.conn == nil {
		return networking.ErrConnectionNotEstablished
	}

	err := client.conn.Close()
	client.conn = nil
	return err
}
