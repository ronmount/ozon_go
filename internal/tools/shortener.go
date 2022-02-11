package tools

import (
	gonanoid "github.com/matoous/go-nanoid"
)

func GenerateToken() (string, error) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	token, err := gonanoid.Generate(alphabet, 10)
	return token, err
}
