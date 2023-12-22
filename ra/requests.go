package ra

import (
	"bytes"
	"fmt"
	"net/http"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)" +
	" AppleWebKit/537.36 (KHTML, like Gecko)" +
	" Chrome/113.0.0.0 Safari/537.36"

func getArtistBySlugReq(slug Slug, baseUri string) (*http.Request, error) {
	query := fmt.Sprintf(`{"query":"{\n artist(slug:\"%v\"){\n id\n name\n}\n}\n","variables":{}}`, slug)
	reqBody := []byte(query)

	req, err := http.NewRequest("POST", baseUri+"/graphql", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	return req, err
}
