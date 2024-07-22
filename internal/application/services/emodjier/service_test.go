package emojier

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceItem(t *testing.T) {
	serv := data{}
	require.Equal(t, len(serv.Commands()) > 0, true)
}
