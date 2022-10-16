package dbconnection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDb(t *testing.T) {
	_, err := NewDatabadeConnection().GetDb()
	assert := assert.New(t)
	assert.NoError(err)
}
