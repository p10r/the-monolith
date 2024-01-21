package ra

import (
	"pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestNewSlug(t *testing.T) {
	for _, url := range []string{
		"https://ra.co/dj/crilletamalt",
		"https://ra.co/dj/crilletamalt/past-events",
		"  https://ra.co/dj/crilletamalt",
		"https://ra.co/dj/crilletamalt  ",
		"ra.co/dj/crilletamalt  ",
		"https://ra.co/dj/crilletamalt/",
	} {
		t.Run("deserializes "+url, func(t *testing.T) {
			got, err := domain.NewSlug(url)
			want := domain.RASlug("crilletamalt")

			expect.NoErr(t, err)
			expect.Equal(t, got, want)
		})
	}

	t.Run("returns err", func(t *testing.T) {
		for _, url := range []string{
			"",
			"facebook.com/harald",
			"https://ra.co/whatever/crilletamalt/",
		} {
			t.Run("maps"+url+"to user id", func(t *testing.T) {
				_, err := domain.NewSlug(url)

				expect.Err(t, err)
			})
		}
	})

}
