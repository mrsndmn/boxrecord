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

	updateOps []tnt.Operator
}

type BoxTest1IndexedFields struct {
	f1 uint32
	f2 uint32
	f3 uint32
	f4 uint32
}

type BoxTest1PK struct {
	f1 uint32
	f2 uint32
}

func (bt *BoxTest1) SetF3(newF3 uint32) {
	if bt.f3 == newF3 {
		return
	}
	bt.f3 = newF3
	bt.updateOps = append(bt.updateOps, tnt.OpSet(2, tnt.PackInt(bt.f3))) // todo а что будет, если тут будет несколько операций с одним и тем же полем?

	return
}

func (bt *BoxTest1) SetF4(newF4 uint32) {
	if bt.f4 == newF4 {
		return
	}
	bt.f4 = newF4
	bt.updateOps = append(bt.updateOps, tnt.OpSet(3, tnt.PackInt(bt.f4)))

	return
}

func (bt *BoxTest1) Update(ctx context.Context, conn *tnt.Connection) error {
	if len(bt.updateOps) == 0 {
		return nil
	}

	_, err := conn.Exec(ctx, &tnt.Update{
		Tuple: tnt.Tuple{
			tnt.PackInt(bt.f1), // primary index should be enough
			tnt.PackInt(bt.f2), // primary index should be enough
		},
		Ops: bt.updateOps,
	})

	bt.updateOps = nil
	return err
}

func Create(ctx context.Context, conn *tnt.Connection, tupleFields *BoxTest1IndexedFields) (*BoxTest1, error) {
	tuple := tnt.Tuple{
		tnt.PackInt(tupleFields.f1),
		tnt.PackInt(tupleFields.f2),
		tnt.PackInt(tupleFields.f3),
		tnt.PackInt(tupleFields.f4),
	}
	_, err := conn.Exec(ctx, &tnt.Insert{
		Tuple: tuple,
	})
	if err != nil {
		return nil, err
	}

	return &BoxTest1{tupleFields.f1, tupleFields.f2, tupleFields.f3, tupleFields.f4}, nil
}

func (bt *BoxTest1) Delete(ctx context.Context, conn *tnt.Connection, tupleFields *BoxTest1IndexedFields) error {
	_, err := conn.Exec(ctx, &tnt.Delete{
		Tuple: tnt.Tuple{
			tnt.PackInt(bt.f1), // primary index should be enough
			tnt.PackInt(bt.f2), // primary index should be enough
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
