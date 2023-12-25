package expect

import (
	"reflect"
	"testing"
)

//goland:noinspection GoUnusedExportedFunction
func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

//goland:noinspection GoUnusedExportedFunction
func SliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got %v, want %v", got, want)
		return
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
			return
		}
	}
}

//goland:noinspection GoUnusedExportedFunction
func NotEmpty[T comparable](t *testing.T, got []T) {
	t.Helper()
	if len(got) > 0 {
		return
	} else {
		t.Errorf("got %v, want at least 1", got)
	}
}

//goland:noinspection GoUnusedExportedFunction
func DeepEqual[T any](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

//goland:noinspection GoUnusedExportedFunction
func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("didn't want %v", got)
	}
}

//goland:noinspection GoUnusedExportedFunction
func Len[T any](t *testing.T, got []T, want int) {
	if len(got) != want {
		t.Errorf("got length %d, want %d", len(got), want)
	}
}

//goland:noinspection GoUnusedExportedFunction
func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Error("got false, want true")
	}
}

//goland:noinspection GoUnusedExportedFunction
func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Error("got true, want false")
	}
}

//goland:noinspection GoUnusedExportedFunction
func NoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

//goland:noinspection GoUnusedExportedFunction
func Err(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("wanted error")
	}
}
