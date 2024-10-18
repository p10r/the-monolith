package giftbox_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/pedro/giftbox"
	"testing"
)

func TestGift(t *testing.T) {
	t.Run("creates gift", func(t *testing.T) {
		expected := giftbox.Gift{
			ID:       "1",
			Type:     giftbox.TypeImage,
			ImageUrl: "example.com",
		}
		gift, err := giftbox.NewGift("1", giftbox.TypeImage, false, "example.com")

		assert.NoError(t, err)
		assert.Equal(t, expected, gift)
	})

	t.Run("returns err if creating an image gift without image url", func(t *testing.T) {
		_, err := giftbox.NewGift("1", giftbox.TypeImage, false, "")

		assert.Error(t, err)
	})

	t.Run("returns err if image url is set for a non-image gift", func(t *testing.T) {
		_, err := giftbox.NewGift("1", giftbox.TypeWish, false, "example.com")

		assert.Error(t, err)
	})

}
