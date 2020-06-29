package tntrecord

import (
	"context"
	"log"
	"testing"

	"github.com/lomik/go-tnt"
	"github.com/stretchr/testify/require"
)

func TestBox(t *testing.T) {
	connector := tnt.New("localhost:10000", &tnt.Options{})
	conn, err := connector.Connect()
	require.NoError(t, err)

	defer conn.Close()

	res, err := conn.Exec(context.TODO(), &tnt.Insert{
		Tuple: tnt.Tuple{tnt.PackInt(1), tnt.PackInt(1)},
	})
	require.NoError(t, err)
	log.Printf("%x", res)

	res, err = conn.Exec(context.TODO(), &tnt.Insert{
		Tuple: tnt.Tuple{tnt.PackInt(2), tnt.PackInt(2)},
	})
	require.NoError(t, err)
	log.Printf("%x", res)

	t1 := tnt.Tuple{tnt.PackInt(1), tnt.PackInt(1)}
	t2 := tnt.Tuple{tnt.PackInt(2), tnt.PackInt(2)}
	res, err = conn.Exec(context.TODO(), &tnt.Select{
		Index:  0,
		Space:  0,
		Tuples: []tnt.Tuple{t1, t2},
		Limit:  2,
	})
	require.NoError(t, err)

	log.Printf("selected: %x", res)

}
