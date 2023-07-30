package domain

type Artist struct {
	Id   Id
	Name string
}

type Artists []Artist

type NewArtist struct {
	Name string
}

type Id int
