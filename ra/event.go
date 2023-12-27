package ra

type Event struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	StartTime  string `json:"startTime"`
	ContentUrl string `json:"contentUrl"`
	Images     []struct {
		Filename string `json:"filename"`
	} `json:"images"`
	Venue struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		ContentUrl string `json:"contentUrl"`
		Area       struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"area"`
	} `json:"venue"`
}
