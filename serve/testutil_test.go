package serve

import (
	"encoding/json"
	"os"
	"testing"
)

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func matchesFromTestdata(t *testing.T) Matches {
	t.Helper()
	var fixture flashscoreResponse
	if err := json.Unmarshal(mustReadFile(t, "testdata/flashscore-res.json"), &fixture); err != nil {
		t.Fatalf("unmarshal flashscore response: %v", err)
	}
	return fixture.toMatches()
}
