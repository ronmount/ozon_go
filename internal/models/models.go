package models

type Database interface {
	AddURL(fullLink string) (interface{}, error)
	GetURL(shortLink string) (interface{}, error)
}

type Link struct {
	FullLink  string `json:"full_link"`
	ShortLink string `json:"short_link"`
}
