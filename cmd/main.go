package main

import (
	"errors"
	"github.com/ronmount/ozon_go/internal/backend_server"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	PsqlError      = "psql connection fault"
	WrongStorage   = "wrong storage type"
	MissingStorage = "missing storage type"
	RunError       = "run error"
)

var log = logrus.New()

func main() {
	var info string
	if server, err := backend_server.NewServer(); err != nil {
		if errors.As(err, &models.PSQLStorage{}) {
			info = PsqlError
		} else if errors.As(err, &models.WrongStorageType{}) {
			info = WrongStorage
		} else if errors.As(err, &models.MissingStorageType{}) {
			info = MissingStorage
		}
		log.Fatal(info)
	} else if err := server.Run(); err != nil {
		log.Fatal(RunError)
	}
}
