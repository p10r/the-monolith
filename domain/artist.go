package domain

type Artist struct {
	Id   int64
	RAId int64
	Name string
}

type Artists []Artist

type NewArtist struct {
	Name string
}
