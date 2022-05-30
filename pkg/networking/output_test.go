package networking

import "testing"

func TestBytes(t *testing.T) {
	inputs := [][]byte{
		nil,
		{0x00, 0x01, 0x02, 0x03},
	}
	expectedValues := [][]byte{
		{},
		{0x00, 0x01, 0x02, 0x03},
	}

	var out Output
	var res []byte

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteBytes(inputs[i])

		res = out.Bytes()

		if !BytesEqual(res, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], res)
		}
	}
}

func TestSequentialWriteSingleByte(t *testing.T) {
	inputs := []byte{
		0x00,
		0x01,
		0x02,
		0x03,
	}
	expectedValue := inputs

	out := NewOutput()

	for _, input := range inputs {
		out.WriteSingleByte(input)
	}

	if !BytesEqual(out.buf, expectedValue) {
		t.Errorf("Expected %v got %v.", expectedValue, out.buf)
	}
}

func TestSequentialWriteBytes(t *testing.T) {
	inputs := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		{0x00, 0x01, 0x02, 0x03},
	}
	expectedValue := append(inputs[0], inputs[1]...)

	out := NewOutput()

	for _, input := range inputs {
		out.WriteBytes(input)
	}

	if !BytesEqual(out.buf, expectedValue) {
		t.Errorf("Expected %v got %v.", expectedValue, out.buf)
	}
}

func TestWriterInterface(t *testing.T) {
	inputs := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		nil,
	}
	expectedValues := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		{},
	}
	expectedNs := []int{
		4,
		0,
	}
	expectedErrors := []bool{
		false,
		false,
	}

	var out Output
	var n int
	var err error

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		n, err = out.Write(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
		if n != expectedNs[i] {
			t.Errorf("N %d: Expected %v got %v.", i, expectedNs[i], n)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestWriteSingleByte(t *testing.T) {
	inputs := []byte{
		0xAD,
	}
	expectedValues := [][]byte{
		{0xAD},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteSingleByte(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestByteWriterInterface(t *testing.T) {
	inputs := []byte{
		0xAD,
	}
	expectedValues := [][]byte{
		{0xAD},
	}
	expectedErrors := []bool{false}

	var out Output
	var err error

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteByte(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
		if (err != nil) != expectedErrors[i] {
			t.Errorf("Error %d: Expected %v got %v.", i, expectedErrors[i], err != nil)
		}
	}
}

func TestWriteBytes(t *testing.T) {
	inputs := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		nil,
	}
	expectedValues := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		{},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteBytes(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteBigEndianInt16(t *testing.T) {
	inputs := []uint16{
		31333,
	}
	expectedValues := [][]byte{
		{0x7A, 0x65},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteBigEndianInt16(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteLittleEndianInt16(t *testing.T) {
	inputs := []uint16{
		25978,
	}
	expectedValues := [][]byte{
		{0x7A, 0x65},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteLittleEndianInt16(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteBigEndianInt32(t *testing.T) {
	inputs := []uint32{
		1635411314,
	}
	expectedValues := [][]byte{
		{0x61, 0x7A, 0x65, 0x72},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteBigEndianInt32(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteLittleEndianInt32(t *testing.T) {
	inputs := []uint32{
		1919253089,
	}
	expectedValues := [][]byte{
		{0x61, 0x7A, 0x65, 0x72},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteLittleEndianInt32(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteBigEndianInt64(t *testing.T) {
	inputs := []uint64{
		7024038111092526354,
	}
	expectedValues := [][]byte{
		{0x61, 0x7A, 0x65, 0x72, 0x74, 0x79, 0xCD, 0x12},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteBigEndianInt64(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteLittleEndianInt64(t *testing.T) {
	inputs := []uint64{
		1354872603950807649,
	}
	expectedValues := [][]byte{
		{0x61, 0x7A, 0x65, 0x72, 0x74, 0x79, 0xCD, 0x12},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteLittleEndianInt64(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteUVarInt(t *testing.T) {
	inputs := []uint64{
		1354872603950807649,
	}
	expectedValues := [][]byte{
		{225, 244, 149, 147, 199, 174, 222, 230, 18},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteUVarInt(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteVarInt(t *testing.T) {
	inputs := []int64{
		1354872603950807649,
	}
	expectedValues := [][]byte{
		{194, 233, 171, 166, 142, 221, 188, 205, 37},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteVarInt(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteNullTerminatedString(t *testing.T) {
	inputs := []string{
		"azerty",
	}
	expectedValues := [][]byte{
		{0x61, 0x7A, 0x65, 0x72, 0x74, 0x79, 0x00},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteNullTerminatedString(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteString(t *testing.T) {
	inputs := []string{
		"azerty",
	}
	expectedValues := [][]byte{
		{0x06, 0x61, 0x7A, 0x65, 0x72, 0x74, 0x79},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteString(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestWriteRaknetString(t *testing.T) {
	inputs := []string{
		"azerty",
	}
	expectedValues := [][]byte{
		{0x00, 0x06, 0x61, 0x7A, 0x65, 0x72, 0x74, 0x79},
	}

	var out Output

	for i := 0; i < len(inputs); i++ {
		out = NewOutput()
		out.WriteRaknetString(inputs[i])

		if !BytesEqual(out.buf, expectedValues[i]) {
			t.Errorf("Value %d: Expected %v got %v.", i, expectedValues[i], out.buf)
		}
	}
}

func TestMergeOutputs(t *testing.T) {
	out1 := NewOutput()
	out2 := NewOutput()

	out1.WriteBigEndianInt32(25)
	out2.WriteString("azerty")

	if !BytesEqual(MergeOutputs(out1, out2).buf, append(out1.buf, out2.buf...)) {
		t.Errorf("Not merged correctly.")
	}
}
