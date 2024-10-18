package giftbox

import (
	"errors"
	"fmt"
)

type GiftType string

const (
	TypeSweet GiftType = "SWEET"
	TypeWish  GiftType = "WISH"
	TypeImage GiftType = "IMAGE"
)

type GiftID string

func (id GiftID) String() string {
	return string(id)
}

type Gift struct {
	ID       GiftID
	Type     GiftType
	Redeemed bool
	// Only set for type "IMAGE". Might be moved to a separate struct later
	ImageUrl string
}

func NewGift(
	id GiftID,
	giftType GiftType,
	redeemed bool,
	imageUrl string,
) (Gift, error) {
	if giftType == TypeImage && imageUrl == "" {
		return Gift{}, errors.New("imageUrl has to be set for gifts of type image")
	}
	if giftType != TypeImage && imageUrl != "" {
		return Gift{}, fmt.Errorf("imageUrl should not be set for %s. ID: %v", giftType, id)
	}

	return Gift{
		ID:       id,
		Type:     giftType,
		Redeemed: redeemed,
		ImageUrl: imageUrl,
	}, nil
}

type Gifts []Gift

func (g Gifts) findByID(reqId string) (Gift, bool) {
	giftsByID := make(map[string]Gift)
	for _, gift := range g {
		giftsByID[gift.ID.String()] = gift
	}
	gift, ok := giftsByID[reqId]
	return gift, ok
}
