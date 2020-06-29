package tntrecord

import (
	"context"
	"testing"

	"github.com/lomik/go-tnt"
	"github.com/stretchr/testify/assert"
)

func TestBox(t *testing.T) {
	connector := tnt.New("localhost:10000", &tnt.Options{})
	conn, err := connector.Connect()
	assert.NoError(t, err)
	defer conn.Close()

	conn.Exec(context.TODO(), &tnt.Insert{})

}
