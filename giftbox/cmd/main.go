package main

import (
	json2 "encoding/json"
	"github.com/p10r/pedro/giftbox"
	"github.com/skip2/go-qrcode"
	"log"
	"net/http"
	"os"
	"time"
)

const baseUrl = "https://pedro-go.fly.dev"

type GiftRef struct {
	url string
	id  string
}

// Generate QR-Codes
// Run:
// direnv allow . && go run main.go
func main() {
	apiKey := os.Getenv("GIFT_BOX_API_KEY")
	if apiKey == "" {
		log.Fatal("GIFT_BOX_API_KEY not set")
	}

	gifts := fetchGifts(apiKey)
	refs := toGiftIDs(gifts)
	writeQRCodes(refs)
}

func writeQRCodes(refs []GiftRef) {
	for _, ref := range refs {
		err := qrcode.WriteFile(ref.url, qrcode.Medium, 256, "qr-"+ref.id+".png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func toGiftIDs(gifts giftbox.Gifts) []GiftRef {
	var refs []GiftRef
	for _, gift := range gifts {
		ref := GiftRef{
			url: baseUrl + "/allGiftsRes?pending-only=true" + string(gift.ID),
			id:  string(gift.ID),
		}
		refs = append(refs, ref)
	}
	return refs
}

func fetchGifts(key string) giftbox.Gifts {
	req, err := http.NewRequest("GET", baseUrl+"/gifts?pending-only=true", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gift-Box-Api-Key", key)

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d", res.StatusCode)
	}
	var allGiftsRes giftbox.AllGiftsRes
	err = json2.NewDecoder(res.Body).Decode(&allGiftsRes)
	if err != nil {
		log.Fatal(err)
	}

	gifts := allGiftsRes.Gifts
	return gifts
}
