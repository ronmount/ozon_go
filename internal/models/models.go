package models

type Database interface {
	AddURL(fullLink string) (interface{}, int)
	GetURL(shortLink string) (interface{}, int)
}

type Link struct {
	FullLink  string `json:"full_link"`
	ShortLink string `json:"short_link"`
}
