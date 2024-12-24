package testutil

import (
	"encoding/json"
	"github.com/alecthomas/assert/v2"
	"testing"
)

func PrettyPrinted(t *testing.T, v any) []byte {
	marshal, err := json.MarshalIndent(v, "", " ")
	assert.NoError(t, err)
	return marshal
}
