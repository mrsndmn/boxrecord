package tntrecord

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/lomik/go-tnt"
	"github.com/stretchr/testify/require"
)

func TestBoxMethod(t *testing.T) {
	connector := tnt.New("localhost:10000", &tnt.Options{})
	conn, err := connector.Connect()
	require.NoError(t, err)

	defer conn.Close()

	res, err := conn.Exec(context.TODO(), &tnt.Insert{
		Tuple: tnt.Tuple{tnt.PackInt(1), tnt.PackInt(1), tnt.PackInt(1)},
	})
	require.NoError(t, err)
	log.Printf("%x", res)

	res, err = conn.Exec(context.TODO(), &tnt.Insert{
		Tuple: tnt.Tuple{tnt.PackInt(2), tnt.PackInt(2), tnt.PackInt(1)},
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

	log.Printf("selected: %x", res[0])

}

// tntboxrecord BoxTest
type BoxTest struct {
	f1 int `box:index1`
	f2 int `box:index1`
	f3 int // any count of other fields
}

func TestTntRecord(t *testing.T) {
	connector := tnt.New("localhost:10000", &tnt.Options{})
	conn, err := connector.Connect()
	require.NoError(t, err)

	defer conn.Close()

	record, err := Create(context.TODO(), conn, &BoxTest1IndexedFields{0, 0, 0, 0})
	require.NoError(t, err)
	fmt.Printf("record %s\n", record)

	record.SetF3(1)
	err = record.Update(context.TODO(), conn)
	require.NoError(t, err)

	selectedRecord, err := SelectByPK(context.TODO(), conn, record.GetPK())
	require.NoError(t, err)
	fmt.Printf("selected %s\n", selectedRecord)
	fmt.Printf("updated  %s\n", record)
	require.True(t, selectedRecord.Equals(record))

	err = record.Delete(context.TODO(), conn)
	require.NoError(t, err)

	selectedRecord, err = SelectByPK(context.TODO(), conn, record.GetPK())
	require.NoError(t, err)
	require.Nil(t, selectedRecord)

	record1, err := Create(context.TODO(), conn, &BoxTest1IndexedFields{0, 1, 0, 0})
	require.NoError(t, err)
	fmt.Printf("record1 %s\n", record1)

	record2, err := Create(context.TODO(), conn, &BoxTest1IndexedFields{0, 2, 0, 0})
	require.NoError(t, err)
	fmt.Printf("record2 %s\n", record2)

	records, err := SelectMultiByPK(context.TODO(), conn, []*BoxTest1PK{record1.GetPK(), record2.GetPK()})
	require.Len(t, records, 2)

	require.True(t, record1.Equals(records[0]))
	require.True(t, record2.Equals(records[1]))
}
