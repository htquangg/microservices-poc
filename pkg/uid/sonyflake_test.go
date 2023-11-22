package uid

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSonyflake(t *testing.T) {
	sf := NewSonyflake()

	nextID := sf.ID()
	assert.NotNil(t, nextID)
}

func TestSonyflake_WithServerIDEqualsZero(t *testing.T) {
	sf1 := NewSonyflake()
	assert.Zero(t, sf1.ServerID())

	var nextID string

	nextID = sf1.ID()
	assert.NotNil(t, nextID)

	os.Setenv("SERVER_ID", "0")

	sf2 := NewSonyflake()
	assert.Zero(t, sf2.ServerID())

	nextID = sf2.ID()
	assert.NotNil(t, nextID)
}

func TestSonyflake_WithServerID(t *testing.T) {
	os.Setenv("SERVER_ID", "100")

	sf := NewSonyflake()
	assert.NotZero(t, sf.ServerID())

	nextID := sf.ID()
	assert.NotNil(t, nextID)
}
