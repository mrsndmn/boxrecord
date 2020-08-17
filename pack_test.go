package tntrecord

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func (tup *IProtoTestStruct) Size() int {
	return 8 + 4 + 8 + 4 + len(tup.bytes)
}

func (tup *IProtoTestStruct) IProtoPack(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, tup.int64key)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.int32key)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.uint64key)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, tup.uint32key)
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
	buf := make([]byte, tup.Size())

	binary.LittleEndian.PutUint64(buf[0:8], uint64(tup.int64key))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(tup.int32key))
	binary.LittleEndian.PutUint64(buf[12:20], tup.uint64key)
	binary.LittleEndian.PutUint32(buf[20:24], tup.uint32key)
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

	if  buf.Len() != tup.Size() {
		t.Errorf("Tuple size and written size are not equal. Bub len %d. Tuple size: %d", buf.Len(), tup.Size())
	}
}
