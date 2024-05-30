package ra_test

import (
	"errors"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/pedro/ra"
	"net/http/httptest"
	"testing"
)

func TestArtist(t *testing.T) {
	t.Run("deserializes success response", func(t *testing.T) {
		res := GraphQLRes(`
			{
			    "data": {
			        "artist": {
			            "id": "106972",
			            "name": "Sinamin"
			        }
			    }
			}`)

		got, err := ra.NewArtist(res.Result())
		want := ra.Artist{RAID: "106972", Name: "Sinamin"}

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("handles not found response", func(t *testing.T) {
		res := GraphQLRes(
			`{
				    "data": {
				        "artist": null
				    }
				}`,
		)

		_, err := ra.NewArtist(res.Result())

		expect.Err(t, err)
		expect.True(t, errors.Is(err, ra.ErrSlugNotFound))
	})

	t.Run("reports non-200 response", func(t *testing.T) {
		res := httptest.NewRecorder()
		res.WriteHeader(500)

		_, err := ra.NewArtist(res.Result())

		expect.Err(t, err)
	})

	t.Run("reports invalid json", func(t *testing.T) {
		res := GraphQLRes(`{"missingHyphen: ""}`)

		_, err := ra.NewArtist(res.Result())

		want := errors.New(`JSON deserialization error. Body: {"missingHyphen: ""}`)
		expect.Err(t, err)
		expect.Equal(t, err.Error(), want.Error())
	})
}

func GraphQLRes(body string) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	res.WriteHeader(200)
	res.Header().Add("Content-Type", "application/graphql-response+json; charset=utf-8")
	_, _ = res.Write([]byte(body))
	return res
}
