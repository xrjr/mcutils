// networking provides all network abstractions : connection read and write, parsing some data types, errors, etc...
package networking

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	MaximumUDPDatagramLength int = 65535
)

// Usual networking errors common to clients
var (
	ErrConnectionNotEstablished     error = errors.New("connection hasn't been established yet. Call Connect method to establish connection")
	ErrConnectionAlreadyEstablished error = errors.New("connection has already been established. If you want to reopen a connection for this client, you have to call Disconnect first")
)

// Conn is common interface between TCP and UDP connections.
type Conn interface {
	ExecuteRequest(Output) (Input, error)
	Close() error
}

// TCPConn is a tcp connection.
type TCPConn struct {
	conn *net.TCPConn
}

// DialTCPOptions are the options for the DialTCP function.
// An empty struct (all fields set to false) is considered as the default behavior for the DialTCP function.
type DialTCPOptions struct {
	SkipSRVLookup bool
	DialTimeout   time.Duration
}

// DialTCP resolve TCP address and connects to the address using TCP.
func DialTCP(hostname string, port int, options DialTCPOptions) (*TCPConn, error) {
	var _hostname string = hostname
	var _port int = port

	if !options.SkipSRVLookup {
		_, addrs, err := net.LookupSRV("minecraft", "tcp", hostname)
		if err == nil && len(addrs) > 0 {
			_hostname = addrs[0].Target
			_port = int(addrs[0].Port)
		}
	}

	c, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", _hostname, _port), options.DialTimeout)
	if err != nil {
		return nil, err
	}

	return &TCPConn{
		conn: c.(*net.TCPConn),
	}, nil
}

// Send sends output to the connection, waits for response and returns the connection input.
// For TCP connections, as they can be read in multiple time, the connection is simply passed as the reader of the response.
func (tcpc TCPConn) Send(req Output) (Input, error) {
	if tcpc.conn == nil {
		return Input{}, ErrConnectionNotEstablished
	}
	_, err := tcpc.conn.Write(req.buf)
	if err != nil {
		return Input{}, err
	}

	return NewInput(tcpc.conn), nil
}

// SetReadDeadline sets the read deadline of the underlying connection.
func (tcpc TCPConn) SetReadDeadline(d time.Duration) error {
	if tcpc.conn == nil {
		return ErrConnectionNotEstablished
	}
	return tcpc.conn.SetReadDeadline(time.Now().Add(d))
}

// Close closes the connection.
func (tcpc TCPConn) Close() error {
	if tcpc.conn == nil {
		return ErrConnectionNotEstablished
	}
	return tcpc.conn.Close()
}

// UDPConn is a udp connection.
type UDPConn struct {
	conn *net.UDPConn
}

// DialUDPOptions are the options for the DialUDP function.
// An empty struct (all fields set to false) is considered as the default behavior for the DialUDP function.
type DialUDPOptions struct {
	SkipSRVLookup                bool
	ForceUDPProtocolForSRVLookup bool
	DialTimeout                  time.Duration
}

// DialTCP resolve UDP address and connects to the address using UDP.
func DialUDP(hostname string, port int, options DialUDPOptions) (*UDPConn, error) {
	var _hostname string = hostname
	var _port int = port
	var SRVLookupProtocol string = "tcp"

	if options.ForceUDPProtocolForSRVLookup {
		SRVLookupProtocol = "udp"
	}

	if !options.SkipSRVLookup {
		_, addrs, err := net.LookupSRV("minecraft", SRVLookupProtocol, hostname)
		if err == nil && len(addrs) > 0 {
			_hostname = addrs[0].Target
			_port = int(addrs[0].Port)
		}
	}

	c, err := net.DialTimeout("udp4", fmt.Sprintf("%s:%d", _hostname, _port), options.DialTimeout)
	if err != nil {
		return nil, err
	}

	return &UDPConn{
		conn: c.(*net.UDPConn),
	}, nil
}

// Send sends output to the connection, waits for response and returns the connection input.
// For UDP connections, as they cannot be read in multiple time, the connection is read a single time and loaded into a buffer of size MaximumUDPDatagramLength.
// UDP datagram length should not be over MaximumUDPDatagramLength, so the entire datagram should be loaded. A *bytes.Buffer is the passed as the reader for the response.
func (udpc UDPConn) Send(out Output) (Input, error) {
	if udpc.conn == nil {
		return Input{}, ErrConnectionNotEstablished
	}
	_, err := udpc.conn.Write(out.buf)
	if err != nil {
		return Input{}, err
	}

	var buf [MaximumUDPDatagramLength]byte
	n, err := udpc.conn.Read(buf[:])
	if err != nil {
		return Input{}, err
	}

	return NewInput(bytes.NewBuffer(buf[:n])), nil
}

// SetReadDeadline sets the read deadline of the underlying connection.
func (udpc UDPConn) SetReadDeadline(d time.Duration) error {
	if udpc.conn == nil {
		return ErrConnectionNotEstablished
	}
	return udpc.conn.SetReadDeadline(time.Now().Add(d))
}

// Close closes the connection.
func (udpc UDPConn) Close() error {
	if udpc.conn == nil {
		return ErrConnectionNotEstablished
	}
	return udpc.conn.Close()
}
