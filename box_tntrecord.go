package tntrecord

import (
	"context"
	"fmt"

	"github.com/lomik/go-tnt" //  почему-то мутаторы тут не поддерживаются..
)

type BoxTest1 struct {
	f1 uint32 `box:index1`
	f2 uint32 `box:index1`
	f3 uint32 `box:index2`

	f4 uint32 // any other field
}

type BoxTest1PK struct {
	f1 uint32
	f2 uint32
}

type BoxTest1IndexedFields struct {
	f1 uint32
	f2 uint32
	f3 uint32
}

func Create(pk *BoxTest1IndexedFields) *BoxTest1 {

	return nil
}

func (bt *BoxTest1) Insert(ctx context.Context, conn *tnt.Connection) error {
	_, err := conn.Exec(ctx, &tnt.Insert{
		Tuple: tnt.Tuple{tnt.PackInt(2), tnt.PackInt(2)},
	})
	return err
}

func (bt *BoxTest1) Update(ctx context.Context, conn *tnt.Connection) error {
	_, err := conn.Exec(ctx, &tnt.Update{
		Tuple: tnt.Tuple{tnt.PackInt(2), tnt.PackInt(2)},
		Ops: []tnt.Operator{
			tnt.OpSet(0, tnt.PackInt(bt.f1)),
			tnt.OpSet(1, tnt.PackInt(bt.f2)),
			tnt.OpSet(2, tnt.PackInt(bt.f3)),
		},
	})
	return err
}

func SelectByf1f2(ctx context.Context, conn *tnt.Connection, pk *BoxTest1PK) (*BoxTest1, error) {
	idxTuple := tnt.Tuple{tnt.PackInt(pk.f1), tnt.PackInt(pk.f2)}
	res, err := conn.Exec(ctx, &tnt.Select{
		Index:  0,
		Tuples: []tnt.Tuple{idxTuple},
	})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	// для нулевых индексов и коробок без шардирования
	// todo это же правильно?
	if len(res) > 1 {
		panic("First index for box BoxTest1 expected to be uniq!")
	}

	selectedTuple := res[0]
	record := &BoxTest1{}

	err = parseBoxTest1Tuple(selectedTuple, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func SelectMultiByf1f2(ctx context.Context, conn *tnt.Connection, pks []*BoxTest1PK) ([]*BoxTest1, error) {

	if len(pks) == 0 {
		return nil, nil
	}

	selectTuples := make([]tnt.Tuple, len(pks))
	for i, pk := range pks {
		selectTuples[i] = append(selectTuples[i], tnt.PackInt(pk.f1), tnt.PackInt(pk.f2))
	}

	selectRes, err := conn.Exec(ctx, &tnt.Select{
		Index:  0,
		Tuples: selectTuples,
	})

	if err != nil {
		return nil, err
	}

	if len(selectRes) == 0 {
		return nil, nil
	}

	records := make([]BoxTest1, len(selectRes))
	recordsPtrs := make([]*BoxTest1, 0, len(selectRes))
	for i, selectedTuple := range selectRes {

		record := records[i]
		err := parseBoxTest1Tuple(selectedTuple, &record)
		if err != nil {
			// todo или стоит вернуть, что заселектил, но ошибку тоже отдать?
			return nil, err
		}

		recordsPtrs = append(recordsPtrs, &record)
	}

	return recordsPtrs, nil
}

func parseBoxTest1Tuple(selectedTuple tnt.Tuple, record *BoxTest1) error {

	// unpack index fields
	if len(selectedTuple) < 2 {
		return fmt.Errorf("Selected tuple cardinality mismatch. Expected at least 2 elements in tuple. Got: %x", selectedTuple)
	}
	if len(selectedTuple[0]) < 4 {
		return fmt.Errorf("BoxTest1: index field 0 expected to be uint32 with size of 4 bytes. Got: %x", selectedTuple[0])
	}
	if len(selectedTuple[1]) < 4 {
		return fmt.Errorf("BoxTest1: index field 1 expected to be uint32 with size of 4 bytes. Got: %x", selectedTuple[1])
	}

	record.f1 = tnt.UnpackInt(selectedTuple[0])
	record.f2 = tnt.UnpackInt(selectedTuple[1])

	// unpack another fields only if they are in box
	for i := 2; i < len(selectedTuple); i++ {
		switch i {
		case 2:
			if len(selectedTuple[3]) < 4 {
				return fmt.Errorf("BoxTest1: field 2 expected to be uint32 with size of 4 bytes. Got: %x", selectedTuple[3])
			}
			record.f3 = tnt.UnpackInt(selectedTuple[3])
		// case 3: ...
		default:
			// found new field that was not described by golang struct
			break
		}
	}
	return nil
}

