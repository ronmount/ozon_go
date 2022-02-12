package database

import (
	"github.com/jackc/pgx"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/ronmount/ozon_go/internal/my_errors"
	"github.com/ronmount/ozon_go/internal/tools"
)

const (
	selectShortByFull = `
		SELECT short_link 
		FROM links 
		WHERE full_link = $1
	`
	selectFullByShort = `
		SELECT full_link 
		FROM links 
		WHERE short_link = $1
	`
	insertNewLink = `
		INSERT INTO links 
		(full_link, short_link) 
		VALUES ($1, $2)
	`
)

type PSQLStorage struct {
	pool *pgx.ConnPool
}

func NewPSQLStorage(config pgx.ConnPoolConfig) (*PSQLStorage, error) {
	if pool, err := pgx.NewConnPool(config); err == nil {
		return &PSQLStorage{pool}, nil
	} else {
		return nil, err
	}
}

func (dbs *PSQLStorage) checkLinkAlreadyExists(fullLink string) (interface{}, error) {
	conn, err := dbs.pool.Acquire()
	if err != nil {
		return nil, err
	}
	defer dbs.pool.Release(conn)

	var shortLink string
	err = conn.QueryRow(selectShortByFull, fullLink).Scan(&shortLink)
	if err != nil {
		return nil, nil
	} else {
		return models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	}
}

func (dbs *PSQLStorage) AddURL(fullLink string) (interface{}, error) {

	if link, err := dbs.checkLinkAlreadyExists(fullLink); err != nil {
		return nil, my_errors.HTTP500{}
	} else if link != nil {
		return link, nil
	}

	conn, err := dbs.pool.Acquire()
	if err != nil {
		return nil, my_errors.HTTP500{}
	}
	defer dbs.pool.Release(conn)

	if shortLink, err := tools.GenerateToken(); err != nil {
		return nil, my_errors.HTTP500{}
	} else if _, err = conn.Exec(insertNewLink, fullLink, shortLink); err != nil {
		return nil, my_errors.HTTP500{}
	} else {
		return models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	}
}

func (dbs *PSQLStorage) GetURL(shortLink string) (interface{}, error) {
	conn, err := dbs.pool.Acquire()
	if err != nil {
		return nil, my_errors.HTTP500{}
	}
	defer dbs.pool.Release(conn)

	var fullLink string
	err = conn.QueryRow(selectFullByShort, shortLink).Scan(&fullLink)
	if err != nil {
		return nil, my_errors.HTTP404{}
	} else {
		return models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	}
}
