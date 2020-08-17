package tntrecord

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"io/ioutil"
)

type IProtoTestStruct struct {
	int64key  int64
	int32key  int32
	uint64key uint64
	uint32key uint32
	bytes     []byte
}

const tupleFieldLength = 4

func (tup *IProtoTestStruct) IProtoSize() int {
	return tupleFieldLength * 5 + (8 + 4 + 8 + 4 + len(tup.bytes))
}

func (tup *IProtoTestStruct) IProtoPack(w io.Writer) error {

	err := binary.Write(w, binary.LittleEndian, uint32(8))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.int64key)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, uint32(4))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.int32key)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, uint32(8))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.uint64key)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, uint32(4))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.uint32key)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(tup.bytes)))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.bytes)
	if err != nil {
		return err
	}

	return nil
}

func (tup *IProtoTestStruct) IProtoBulkPack(w io.Writer) error {
	buf := make([]byte, tup.IProtoSize())

	// https://github.com/tarantool/tarantool/blob/stable/doc/box-protocol.txt#L220
	binary.LittleEndian.PutUint32(buf[0:4], 4)
	binary.LittleEndian.PutUint64(buf[4:12], uint64(tup.int64key))
	binary.LittleEndian.PutUint32(buf[12:16], 4)
	binary.LittleEndian.PutUint32(buf[16:20], uint32(tup.int32key))
	binary.LittleEndian.PutUint32(buf[20:24], 4)
	binary.LittleEndian.PutUint64(buf[24:32], tup.uint64key)
	binary.LittleEndian.PutUint32(buf[32:36], 4)
	binary.LittleEndian.PutUint32(buf[36:40], tup.uint32key)
	binary.LittleEndian.PutUint32(buf[40:44], uint32(len(tup.bytes)))
	copy(buf[24:(24+len(tup.bytes))], tup.bytes[:])
	_, err := w.Write(buf)

	return err
}

func BenchmarkIProtoPack(b *testing.B) {
	tup := IProtoTestStruct{
		int64key:  1,
		int32key:  1,
		uint32key: 1,
		uint64key: 1,
		bytes:     []byte{0x0, 0x1},
	}

	for n := 0; n < b.N; n++ {
		tup.IProtoPack(ioutil.Discard)
	}
}

func BenchmarkIProtoBulkPack(b *testing.B) {
	tup := IProtoTestStruct{
		int64key:  1,
		int32key:  1,
		uint32key: 1,
		uint64key: 1,
		bytes:     []byte{0x0, 0x1},
	}

	for n := 0; n < b.N; n++ {
		tup.IProtoBulkPack(ioutil.Discard)
	}
}

func TestIProtoBulkPack(t *testing.T) {
	tup := IProtoTestStruct{
		int64key:  1,
		int32key:  1,
		uint32key: 1,
		uint64key: 1,
		bytes:     []byte{0x1, 0x0, 0x1},
	}

	buf := new(bytes.Buffer)
	tup.IProtoBulkPack(buf)

	if buf.Len() != tup.IProtoSize() {
		t.Errorf("Tuple size and written size are not equal. Bub len %d. Tuple size: %d", buf.Len(), tup.IProtoSize())
	}
}
