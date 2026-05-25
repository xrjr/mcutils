package networking

import "encoding/binary"

const (
	SEGMENT_BITS byte = 0x7F
	CONTINUE_BIT byte = 0x80
)

// Output represents a connection output (i.e. what's written to the connection). It wraps several helpers to write to this output.
type Output struct {
	buf []byte
}

// NewOutput returns a well-formed Output.
func NewOutput() Output {
	return Output{
		buf: make([]byte, 0),
	}
}

// Bytes returns the underlying buffer.
func (out *Output) Bytes() []byte {
	return out.buf
}

// Write is just a wrapper around WriteBytes to make Request implement io.Writer
func (out *Output) Write(buf []byte) (int, error) {
	out.WriteBytes(buf)
	return len(buf), nil
}

// WriteByte is equivalent to WriteSingleByte but implements io.ByteWriter interface
func (out *Output) WriteByte(b byte) error {
	out.WriteSingleByte(b)
	return nil
}

// WriteSingleByte writes a single byte to the output.
func (out *Output) WriteSingleByte(b byte) {
	out.buf = append(out.buf, b)
}

// WriteBytes writes a slice of bytes to the output.
func (out *Output) WriteBytes(b []byte) {
	out.buf = append(out.buf, b...)
}

// WriteBigEndianInt16 writes a big endian 2-bytes int (short) to the output.
func (out *Output) WriteBigEndianInt16(i uint16) {
	int16Buf := make([]byte, 2)
	binary.BigEndian.PutUint16(int16Buf, i)
	out.WriteBytes(int16Buf)
}

// WriteLittleEndianInt16 writes a little endian 2-bytes int (short) to the output.
func (out *Output) WriteLittleEndianInt16(i uint16) {
	int16Buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(int16Buf, i)
	out.WriteBytes(int16Buf)
}

// WriteBigEndianInt32 writes a big endian 4-bytes int to the output.
func (out *Output) WriteBigEndianInt32(i uint32) {
	int32Buf := make([]byte, 4)
	binary.BigEndian.PutUint32(int32Buf, i)
	out.WriteBytes(int32Buf)
}

// WriteLittleEndianInt32 writes a little endian 4-bytes int to the output.
func (out *Output) WriteLittleEndianInt32(i uint32) {
	int32Buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(int32Buf, i)
	out.WriteBytes(int32Buf)
}

// WriteBigEndianInt64 writes a big endian 8-bytes int (long) to the output.
func (out *Output) WriteBigEndianInt64(i uint64) {
	int64Buf := make([]byte, 8)
	binary.BigEndian.PutUint64(int64Buf, i)
	out.WriteBytes(int64Buf)
}

// WriteLittleEndianInt64 writes a little endian 8-bytes int (long) to the output.
func (out *Output) WriteLittleEndianInt64(i uint64) {
	int64Buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(int64Buf, i)
	out.WriteBytes(int64Buf)
}

// WriteVarInt writes a varint to the output.
func (out *Output) WriteVarInt(i int32) {
	v := uint32(i)
	for c := 0; c < 5; c++ {
		currentByte := byte(v & uint32(SEGMENT_BITS))
		v >>= 7
		if v != 0 {
			currentByte |= CONTINUE_BIT
		}
		out.WriteByte(currentByte)

		if v == 0 {
			return
		}
	}
}

// WriteVarLong writes a varlong to the output.
func (out *Output) WriteVarLong(i int64) {
	v := uint64(i)
	for c := 0; c < 10; c++ {
		currentByte := byte(v & uint64(SEGMENT_BITS))
		v >>= 7
		if v != 0 {
			currentByte |= CONTINUE_BIT
		}
		out.WriteByte(currentByte)
		if v == 0 {
			return
		}
	}
}

// WriteNullTerminatedString writes a null terminated string the the output.
func (out *Output) WriteNullTerminatedString(s string) {
	out.WriteBytes([]byte(s))
	out.WriteSingleByte(0)
}

// WriteString writes a standard minecraft protocol string to the output.
// It is a UTF-8 string prefixed with its size in bytes as an unsigned varint.
func (out *Output) WriteString(s string) {
	out.WriteVarInt(int32(len(s)))
	out.WriteBytes([]byte(s))
}

// WriteRaknetString writes a raknet string to the output.
// It is a UTF-8 string prefixed with its size in bytes as a big endian unsigned short.
func (out *Output) WriteRaknetString(s string) {
	out.WriteBigEndianInt16(uint16(len(s)))
	out.WriteBytes([]byte(s))
}

// MergeOutputs merge buffers of two outputs, creating a new output and without modifying any of the merged output buffer.
func MergeOutputs(out1, out2 Output) Output {
	out := NewOutput()
	out.WriteBytes(out1.buf)
	out.WriteBytes(out2.buf)
	return out
}
