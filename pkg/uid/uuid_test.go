package uid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	id := UUIDV4()

	require.NotNil(t, id)
}
