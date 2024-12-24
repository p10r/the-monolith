package testutil

import (
	json2 "encoding/json"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/flashscore"
	"os"
	"testing"
)

func MustReadFile(tb testing.TB, path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		tb.Fatalf("Error trying to load file: %v", err)
	}
	return content
}

func RawFlashscoreRes(tb testing.TB) []byte {
	return MustReadFile(tb, "../testdata/flashscore-res.json")
}

func FlashscoreRes(tb testing.TB) flashscore.Response {
	var res flashscore.Response
	err := json2.Unmarshal(RawFlashscoreRes(tb), &res)
	if err != nil {
		tb.Fatalf("Could not unmarshal raw flashscore response: %v", err)
	}
	return res
}

func UntrackedMatches(tb testing.TB) domain.UntrackedMatches {
	var res flashscore.Response
	err := json2.Unmarshal(RawFlashscoreRes(tb), &res)
	if err != nil {
		tb.Fail()
	}

	return res.ToUntrackedMatches()
}
