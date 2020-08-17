package iproto

import (
	"encoding/binary"
	"io"
)

type ResponseHeader struct {
	Type       int32
	BodyLength int32
	RequestID  int32
	ReturnCode int32
}

const ResponseHeaderSize = 16

func PackResponseHeader(w io.Writer, header ResponseHeader) error {
	var buf [ResponseHeaderSize]byte

	binary.LittleEndian.PutUint32(buf[0:4], uint32(header.Type))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(header.BodyLength))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(header.RequestID))
	binary.LittleEndian.PutUint32(buf[12:16], uint32(header.RequestID))


	_, err := w.Write(buf[:])

	return err
}

// todo support context
func UnpackResponseHeader(r io.Reader) (*ResponseHeader, error) {
	var buf [ResponseHeaderSize]byte
	n := 0

	for {
		readLen, err := r.Read(buf[n:])
		if err != nil {
			return nil, err
		}
		n += readLen
		if n == ResponseHeaderSize {
			break
		}
	}

	reqHeader := &ResponseHeader{
		Type:       int32(binary.LittleEndian.Uint32(buf[0:4])),
		BodyLength: int32(binary.LittleEndian.Uint32(buf[4:8])),
		RequestID:  int32(binary.LittleEndian.Uint32(buf[8:12])),
		ReturnCode: int32(binary.LittleEndian.Uint32(buf[12:16])),
	}

	return reqHeader, nil
}
