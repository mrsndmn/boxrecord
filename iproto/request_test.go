package iproto

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIprotoPing(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:10000")
	require.NoError(t, err)

	req := RequestHeader{
		RequestID: 1,
		Type:      RequestTypePing,
	}

	err = PackRequestHeader(conn, req)
	require.NoError(t, err)
	fmt.Println("req is ok")

	// жесть, iproto запрос возвращает 12 байт и его надо распаковывать в структуру заголовка запроса(((
	respHeader, err := UnpackRequestHeader(conn)
	require.NoError(t, err)
	fmt.Println("resp is ok")

	require.Equal(t, *respHeader, req)
}
