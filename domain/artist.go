package domain

type Artist struct {
	Id   int64
	RAId string
	Name string
}

type Artists []Artist

type NewArtist struct {
	Name string
}
