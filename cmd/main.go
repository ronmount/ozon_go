package main

import (
	"errors"
	"github.com/ronmount/ozon_go/internal/backend_server"
	"github.com/ronmount/ozon_go/internal/my_errors"
	"github.com/sirupsen/logrus"
)

const (
	PsqlError      = "psql connection fault"
	RedisError     = "redis connection fault"
	WrongStorage   = "wrong storage type"
	MissingStorage = "missing storage type"
	RunError       = "run error"
)

var log = logrus.New()

func main() {
	var info string
	if server, err := backend_server.NewServer(); err != nil {
		if errors.As(err, &my_errors.PSQLStorage{}) {
			info = PsqlError
		} else if errors.As(err, &my_errors.RedisStorage{}) {
			info = RedisError
		} else if errors.As(err, &my_errors.WrongStorageType{}) {
			info = WrongStorage
		} else if errors.As(err, &my_errors.MissingStorageType{}) {
			info = MissingStorage
		}
		log.Fatal(info)
	} else if err := server.Run(); err != nil {
		log.Fatal(RunError)
	}
}
