package tntrecord

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var bulkPackPool sync.Pool

func init() {
	bulkPackPool = sync.Pool{}
	bulkPackPool.New = func() interface{} {
		buf := make([]byte, 100)
		return buf
	}
}

func getBulkPackBuffer() []byte {
	return bulkPackPool.Get().([]byte)
}

type IProtoTestStruct struct {
	int64key  int64
	int32key  int32
	uint64key uint64
	uint32key uint32
	bytes     []byte // todo max size 10 bytes
}

const tupleFieldLength = 4
const IProtoTestStructBytesMaxSize = 10
const ExpectedMaxTupleSize = tupleFieldLength*5 + (8 + 4 + 8 + 4 + IProtoTestStructBytesMaxSize)

func (tup *IProtoTestStruct) IProtoSize() int {
	return tupleFieldLength*5 + (8 + 4 + 8 + 4 + len(tup.bytes))
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
	// todo ./pack_test.go:82:13: make([]byte, 100) escapes to heap
	// даже с константным размером этот буфер эскейпится. Думаю, что из-за
	// интерфейса в io.Writer
	// Кейс с пуллом работает помедленнее на бенчмарках, но GC будет проще жить с пуллом..
	// но в целом, можно надрочиться и сделать так, чтобы не было аллокаций

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
	_, err := w.Write(buf[:(24 + len(tup.bytes))])

	return err
}

func (tup *IProtoTestStruct) IProtoBulkPackWithPool(w io.Writer) error {

	buf := getBulkPackBuffer()
	if tup.IProtoSize() > len(buf) {
		fmt.Printf("%d, %d\n", tup.IProtoSize(), len(buf))
		// хардкор! не надо так, надо поддержать увеличение буфера!
		panic("can't pack IProtoTestStruct. Its too big")
	}

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
	_, err := w.Write(buf[:(24 + len(tup.bytes))])

	bulkPackPool.Put(buf)

	return err
}

func (tup *IProtoTestStruct) IProtoBulkPackConn(conn *net.TCPConn) error {
	var buf []byte
	if ExpectedMaxTupleSize < tup.IProtoSize() {
		buf = make([]byte, tup.IProtoSize())
	} else {
		buf = make([]byte, ExpectedMaxTupleSize)
	}

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
	_, err := conn.Write(buf[:(24 + len(tup.bytes))])

	return err
}

func getTestClient(b *testing.B) *net.TCPConn {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	require.NoError(b, err)

	ln, err := net.ListenTCP("tcp", addr)
	require.NoError(b, err)

	// var server net.Conn
	go func() {
		defer ln.Close()
		_, err = ln.Accept()
	}()

	clientAddr, err := net.ResolveTCPAddr(ln.Addr().Network(), ln.Addr().String())
	require.NoError(b, err)

	client, err := net.DialTCP("tcp", nil, clientAddr)
	require.NoError(b, err)

	return client
}

func BenchmarkIProtoBulkPackTCPConn(b *testing.B) {
	tup := IProtoTestStruct{
		int64key:  1,
		int32key:  1,
		uint32key: 1,
		uint64key: 1,
		bytes:     []byte{0x1, 0x0, 0x1},
	}

	client := getTestClient(b)
	var err error

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = tup.IProtoBulkPackConn(client)
		if err != nil {
			b.Error("error IProtoBulkPackConn")
		}
	}
}

func BenchmarkIProtoBulkPackTCPConnUnhappyCase(b *testing.B) {
	tupleBytes := make([]byte, ExpectedMaxTupleSize+1)
	tup := IProtoTestStruct{
		int64key:  1,
		int32key:  1,
		uint32key: 1,
		uint64key: 1,
		bytes:     tupleBytes,
	}

	client := getTestClient(b)
	var err error

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = tup.IProtoBulkPackConn(client)
		if err != nil {
			b.Errorf("error IProtoBulkPackConn %s", err.Error())
		}
	}
}

func BenchmarkIProtoPack(b *testing.B) {
	tup := IProtoTestStruct{
		int64key:  1,
		int32key:  1,
		uint32key: 1,
		uint64key: 1,
		bytes:     []byte{0x0, 0x1},
	}

	client := getTestClient(b)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tup.IProtoPack(client)
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

	client := getTestClient(b)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tup.IProtoBulkPack(client)
	}
}

func BenchmarkIProtoBulkPackWithPool(b *testing.B) {
	tup := IProtoTestStruct{
		int64key:  1,
		int32key:  1,
		uint32key: 1,
		uint64key: 1,
		bytes:     []byte{0x0, 0x1},
	}

	client := getTestClient(b)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tup.IProtoBulkPackWithPool(client)
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
