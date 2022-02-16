package database

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/ronmount/ozon_go/internal/tools"
)

type MemoryStorage struct {
	fullAsKey  map[string]string
	shortAsKey map[string]string
}

func GetMD5Hash(s string) string {
	res := md5.Sum([]byte(s))
	return hex.EncodeToString(res[:])
}

func NewMemoryStorage() (*MemoryStorage, error) {
	rs := &MemoryStorage{}
	rs.shortAsKey = make(map[string]string)
	rs.fullAsKey = make(map[string]string)

	return rs, nil
}

func (ms *MemoryStorage) AddURL(fullLink string) (interface{}, error) {
	var (
		response interface{}
		e        error
	)

	hashFullLink := GetMD5Hash(fullLink)
	if alreadySavedLink, ok := ms.fullAsKey[hashFullLink]; ok {
		response, e = models.Link{FullLink: fullLink, ShortLink: alreadySavedLink}, nil
	} else if shortLink, err := tools.GenerateToken(); err == nil {
		hashShortLink := GetMD5Hash(shortLink)
		ms.fullAsKey[hashFullLink] = shortLink
		ms.shortAsKey[hashShortLink] = fullLink
		response, e = models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	} else {
		response, e = nil, models.HTTP500{}
	}

	return response, e
}

func (ms *MemoryStorage) GetURL(shortLink string) (interface{}, error) {
	var (
		response interface{}
		err      error
	)

	hashShortLink := GetMD5Hash(shortLink)
	if fullLink, ok := ms.shortAsKey[hashShortLink]; ok {
		response, err = models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	} else {
		response, err = nil, models.HTTP404{}
	}

	return response, err
}
