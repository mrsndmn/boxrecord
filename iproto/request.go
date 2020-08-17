package iproto

import (
	"encoding/binary"
	"io"
)

type RequestHeader struct {
	Type       int32
	BodyLength int32
	RequestID  int32
}

const RequstHeaderSize = 12

func PackRequestHeader(w io.Writer, header RequestHeader) error {
	var buf [RequstHeaderSize]byte

	binary.LittleEndian.PutUint32(buf[0:4], uint32(header.Type))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(header.BodyLength))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(header.RequestID))

	_, err := w.Write(buf[:])

	return err
}

// todo support context
func UnpackRequestHeader(r io.Reader) (*RequestHeader, error) {
	var buf [RequstHeaderSize]byte
	n := 0

	for {
		readLen, err := r.Read(buf[n:])
		if err != nil {
			return nil, err
		}
		n += readLen
		if n == RequstHeaderSize {
			break
		}
	}

	reqHeader := &RequestHeader{
		Type: int32(binary.LittleEndian.Uint32(buf[0:4])),
		BodyLength: int32(binary.LittleEndian.Uint32(buf[4:8])),
		RequestID: int32(binary.LittleEndian.Uint32(buf[8:12])),
	}

	return reqHeader, nil
}
