package networking

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Input represents a connection input (i.e. what's read from the connection). It wraps several helpers to read from this input.
type Input struct {
	r io.Reader
}

// NewInput returns a well-formed input.
func NewInput(reader io.Reader) Input {
	return Input{
		r: reader,
	}
}

// Read is just a wrapper around internal reader to make Input implements io.Reader.
func (in *Input) Read(buf []byte) (int, error) {
	return in.r.Read(buf)
}

// ReadByte tries to read a single byte from the input.
// RedByte also implements io.ByteReader interface, which is useful to use binary.ReadUvarint on the Input itself (see method ReadUVarInt).
func (in *Input) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := in.r.Read(buf[:])
	return buf[0], err
}

// ReadBytes tries to read a slice of byte of size n from the input.
func (in *Input) ReadBytes(n int) ([]byte, error) {
	var final *bytes.Buffer = &bytes.Buffer{}
	for i := 0; i < n; i++ {
		b, err := in.ReadByte()
		if err != nil {
			return nil, err
		}
		final.WriteByte(b)
	}
	return final.Bytes(), nil
}

// ReadBigEndianInt16 tries to read a big endian 2-bytes int (short) from the input.
func (in *Input) ReadBigEndianInt16() (uint16, error) {
	buf, err := in.ReadBytes(2)
	return binary.BigEndian.Uint16(buf[0:2]), err
}

// ReadLittleEndianInt16 tries to read a little endian 2-bytes int (short) from the input.
func (in *Input) ReadLittleEndianInt16() (uint16, error) {
	buf, err := in.ReadBytes(2)
	return binary.LittleEndian.Uint16(buf[0:2]), err
}

// ReadBigEndianInt32 tries to read a big endian 4-bytes int from the input.
func (in *Input) ReadBigEndianInt32() (uint32, error) {
	buf, err := in.ReadBytes(4)
	return binary.BigEndian.Uint32(buf[0:4]), err
}

// ReadLittleEndianInt32 tries to read a little endian 4-bytes int from the input.
func (in *Input) ReadLittleEndianInt32() (uint32, error) {
	buf, err := in.ReadBytes(4)
	return binary.LittleEndian.Uint32(buf[0:4]), err
}

// ReadBigEndianInt64 tries to read a big endian 8-bytes int (long) from the input.
func (in *Input) ReadBigEndianInt64() (uint64, error) {
	buf, err := in.ReadBytes(8)
	return binary.BigEndian.Uint64(buf[0:8]), err
}

// ReadLittleEndianInt64 tries to read a little endian 8-bytes int (long) from the input.
func (in *Input) ReadLittleEndianInt64() (uint64, error) {
	buf, err := in.ReadBytes(8)
	return binary.LittleEndian.Uint64(buf[0:8]), err
}

// ReadUVarInt64 tries to read an unsigned varint from the input.
func (in *Input) ReadUVarInt() (uint64, error) {
	tmpUint64, err := binary.ReadUvarint(in)
	if err != nil {
		tmpUint64 = 0
	}
	return tmpUint64, err
}

// ReadVarInt64 tries to read a signed varint from the input.
func (in *Input) ReadVarInt() (int64, error) {
	tmpInt64, err := binary.ReadVarint(in)
	if err != nil {
		tmpInt64 = 0
	}
	return tmpInt64, err
}

// ReadNullTerminatedString tries to read a null terminated string from the input.
func (in *Input) ReadNullTerminatedString() (string, error) {
	var final *bytes.Buffer = &bytes.Buffer{}

	b, err := in.ReadByte()
	if err != nil {
		return "", err
	}

	for b != 0 {
		final.WriteByte(b)
		b, err = in.ReadByte()
		if err != nil {
			return "", err
		}
	}

	return final.String(), nil
}

// ReadString tries to read a standard minecraft protocol string from the input.
// It is a UTF-8 string prefixed with its size in bytes as an unsigned varint.
func (in *Input) ReadString() (string, error) {
	length, err := in.ReadUVarInt()
	if err != nil {
		return "", err
	}
	bytesString, err := in.ReadBytes(int(length))
	return string(bytesString), err
}
