package ra

import (
	"bytes"
	"fmt"
	"net/http"
	"pedro-go/domain"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)" +
	" AppleWebKit/537.36 (KHTML, like Gecko)" +
	" Chrome/113.0.0.0 Safari/537.36"

func newGetArtistReq(slug domain.RASlug, baseUri string) (*http.Request, error) {
	query := fmt.Sprintf(`{
		"query":"{\n artist(slug:\"%v\"){\n id\n name\n}\n}\n",
		"variables":{}
	}`, slug)
	reqBody := []byte(query)

	req, err := http.NewRequest("POST", baseUri+"/graphql", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	return req, err
}

func newGetEventsReq(raId string, start, end time.Time, uri string) (*http.Request, error) {
	fmtStart := start.Format("2006-01-02T03:04:05.000Z")
	fmtEnd := end.Format("2006-01-02T03:04:05.000Z")

	query := fmt.Sprintf(`{
    "operationName": "GET_DEFAULT_EVENTS_LISTING",
    "variables": {
        "indices": [
            "EVENT"
        ],
        "pageSize": 20,
        "page": 1,
        "aggregations": [],
        "filters": [
            {
                "type": "ARTIST",
                "value": "%v"
            },
            {
                "type": "DATERANGE",
                "value": "{\"gte\":\"%v\"}"
            },
            {
                "type": "DATERANGE",
                "value": "{\"lte\":\"%v\"}"
            }
        ],
        "sortOrder": "ASCENDING",
        "sortField": "DATE",
        "baseFilters": [
            {
                "type": "ARTIST",
                "value": "%v"
            },
            {
                "type": "DATERANGE",
                "value": "{\"gte\":\"%v\"}"
            },
            {
                "type": "DATERANGE",
                "value": "{\"lte\":\"%v\"}"
            }
        ]
    },
    "query": "query GET_DEFAULT_EVENTS_LISTING($indices: [IndexType!], $aggregations: [ListingAggregationType!], $filters: [FilterInput], $pageSize: Int, $page: Int, $sortField: FilterSortFieldType, $sortOrder: FilterSortOrderType, $baseFilters: [FilterInput]) {\n listing(indices: $indices, aggregations: [], filters: $filters, pageSize: $pageSize, page: $page, sortField: $sortField, sortOrder: $sortOrder) {\n data {\n ...eventFragment\n __typename\n }\n totalResults\n __typename\n }\n aggregations: listing(indices: $indices, aggregations: $aggregations, filters: $baseFilters, pageSize: 0, sortField: $sortField, sortOrder: $sortOrder) {\n aggregations {\n type\n values {\n value\n name\n __typename\n }\n __typename\n }\n __typename\n }\n}\n\nfragment eventFragment on IListingItem {\n ... on Event {\n id\n title\n attending\n date\n startTime\n contentUrl\n queueItEnabled\n flyerFront\n newEventForm\n images {\n id\n filename\n alt\n type\n crop\n }\n venue {\n id\n name\n contentUrl\n live\n area {\n id\n name\n urlName\n country {\n id\n name\n urlCode\n __typename\n }\n __typename\n }\n }\n pick {\n id\n blurb\n __typename\n }\n __typename\n }\n __typename\n}\n"
}`, raId, fmtStart, fmtEnd, raId, fmtStart, fmtEnd)

	reqBody := []byte(query)

	req, err := http.NewRequest("POST", uri+"/graphql", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	return req, err
}
