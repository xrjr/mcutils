package networking

import (
	"bytes"
	"testing"
)

func BytesEqual(a, b []byte) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestReaderInterface(t *testing.T) {
	buffer := []byte{0x00, 0x01, 0x02}
	expectedValues := [][]byte{
		{0x00, 0x01},
		{0x02},
		{},
	}
	expectedNs := []int{2, 1, 0}
	expectedErrors := []bool{false, false, true}

	decoder := NewInput(bytes.NewBuffer(buffer))

	var n int
	var err error

	for i := 0; i < 2; i++ {
		buf := make([]byte, 2)
		n, err = decoder.Read(buf)

		if !BytesEqual(buf[:n], expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], buf[:n])
		}
		if n != expectedNs[i] {
			t.Errorf("N %d: Expected %v got %v.", i, expectedNs[i], n)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestReadByte(t *testing.T) {
	decoder := NewInput(bytes.NewBuffer([]byte{123, 96}))

	b, err := decoder.ReadByte()

	if b != 123 || err != nil {
		t.Errorf("Expected 123, <nil>. Got %v, %v.", b, err)
	}
}

func TestSequentialReads(t *testing.T) {
	buffer := []byte{0x00, 0x01, 0x02, 0x03}
	expectedValues := []byte{0x00, 0x01, 0x02, 0x03, 0x00}
	expectedErrors := []bool{false, false, false, false, true}

	decoder := NewInput(bytes.NewBuffer(buffer))

	var res byte
	var err error

	for i := 0; i < 5; i++ {
		res, err = decoder.ReadByte()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestReadBytes(t *testing.T) {
	buffer := []byte{0x00, 0x01, 0x02}
	expectedValues := [][]byte{
		{0x00, 0x01},
		nil,
	}
	expectedErrors := []bool{false, true}

	decoder := NewInput(bytes.NewBuffer(buffer))

	var res []byte
	var err error

	for i := 0; i < 2; i++ {
		res, err = decoder.ReadBytes(2)

		if !BytesEqual(res, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestReadBigEndianInt16(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0xA5, 0x96})),
	}
	expectedValues := []uint16{42390}

	var res uint16

	for i := 0; i < len(inputs); i++ {
		res, _ = inputs[i].ReadBigEndianInt16()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
	}
}

func TestReadLittleEndianInt16(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0xA5, 0x96})),
	}
	expectedValues := []uint16{38565}

	var res uint16

	for i := 0; i < len(inputs); i++ {
		res, _ = inputs[i].ReadLittleEndianInt16()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
	}
}

func TestReadBigEndianInt32(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0xA5, 0x96, 0x14, 0x6C})),
	}
	expectedValues := []uint32{2778076268}

	var res uint32

	for i := 0; i < len(inputs); i++ {
		res, _ = inputs[i].ReadBigEndianInt32()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
	}
}

func TestReadLittleEndianInt32(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0xA5, 0x96, 0x14, 0x6C})),
	}
	expectedValues := []uint32{1813288613}

	var res uint32

	for i := 0; i < len(inputs); i++ {
		res, _ = inputs[i].ReadLittleEndianInt32()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
	}
}

func TestReadBigEndianInt64(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0xA5, 0x96, 0x14, 0x6C, 0x69, 0x21, 0xDB, 0x72})),
	}
	expectedValues := []uint64{11931746718617557874}

	var res uint64

	for i := 0; i < len(inputs); i++ {
		res, _ = inputs[i].ReadBigEndianInt64()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
	}
}

func TestReadLittleEndianInt64(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0xA5, 0x96, 0x14, 0x6C, 0x69, 0x21, 0xDB, 0x72})),
	}
	expectedValues := []uint64{8276245476891989669}

	var res uint64
	for i := 0; i < len(inputs); i++ {
		res, _ = inputs[i].ReadLittleEndianInt64()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
	}
}

func TestReadUVarInt(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{232, 201, 171, 166, 15})),
		NewInput(bytes.NewBuffer([]byte{231, 201, 171, 166, 15})),
		NewInput(bytes.NewBuffer([]byte{231, 201, 171, 166})),
	}
	expectedValues := []uint64{4106937576, 4106937575, 0}
	expectedErrors := []bool{false, false, true}

	var res uint64
	var err error

	for i := 0; i < len(inputs); i++ {
		res, err = inputs[i].ReadUVarInt()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestReadVarInt(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{232, 201, 171, 166, 15})),
		NewInput(bytes.NewBuffer([]byte{231, 201, 171, 166, 15})),
		NewInput(bytes.NewBuffer([]byte{231, 201, 171, 166})),
	}
	expectedValues := []int64{2053468788, -2053468788, 0}
	expectedErrors := []bool{false, false, true}

	var res int64
	var err error

	for i := 0; i < len(inputs); i++ {
		res, err = inputs[i].ReadVarInt()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestReadNullTerminatedString(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0x61, 0x7A, 0x65, 0x72, 0x74, 0x79, 0x00, 0x61, 0x7A, 0x65, 0x72, 0x74, 0x79})),
		NewInput(bytes.NewBuffer([]byte{0x61, 0x7A, 0x65, 0x72, 0x74, 0x79})),
		NewInput(bytes.NewBuffer([]byte{})),
	}
	expectedValues := []string{"azerty", "", ""}
	expectedErrors := []bool{false, true, true}

	var res string
	var err error

	for i := 0; i < len(inputs); i++ {
		res, err = inputs[i].ReadNullTerminatedString()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestReadString(t *testing.T) {
	inputs := []Input{
		NewInput(bytes.NewBuffer([]byte{0x06, 0x61, 0x7A, 0x65, 0x72, 0x74, 0x79, 0x61, 0x7A, 0x65, 0x72, 0x74, 0x79})),
		NewInput(bytes.NewBuffer([]byte{0x07, 0x61, 0x7A, 0x65, 0x72, 0x74, 0x79})),
		NewInput(bytes.NewBuffer([]byte{})),
	}
	expectedValues := []string{"azerty", "", ""}
	expectedErrors := []bool{false, true, true}

	var res string
	var err error

	for i := 0; i < len(inputs); i++ {
		res, err = inputs[i].ReadString()

		if res != expectedValues[i] {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}
