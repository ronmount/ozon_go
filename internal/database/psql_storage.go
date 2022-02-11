package database

import (
	"github.com/jackc/pgx"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/ronmount/ozon_go/internal/tools"
	"log"
	"net/http"
)

type PSQLStorage struct {
	pool *pgx.ConnPool
}

func NewPSQLStorage(config pgx.ConnPoolConfig) (*PSQLStorage, error) {
	if pool, err := pgx.NewConnPool(config); err == nil {
		log.Println("Postgres successfully connected.")
		return &PSQLStorage{pool}, nil
	} else {
		return nil, err
	}
}

func (dbs *PSQLStorage) checkLinkAlreadyExists(fullLink string) (interface{}, error) {
	conn, err := dbs.pool.Acquire()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer dbs.pool.Release(conn)
	selectQuery := "select short_link from links where full_link = $1"
	row := conn.QueryRow(selectQuery, fullLink)
	var shortLink string
	err = row.Scan(&shortLink)
	if err != nil {
		return nil, nil
	} else {
		return models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	}
}

func (dbs *PSQLStorage) AddURL(fullLink string) (interface{}, int) {
	link, err := dbs.checkLinkAlreadyExists(fullLink)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}
	if link != nil {
		return link, http.StatusConflict
	}

	conn, err := dbs.pool.Acquire()
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}
	defer dbs.pool.Release(conn)

	shortLink, err := tools.GenerateToken()
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}

	query := "insert into links (full_link, short_link) values ($1, $2)"
	_, err = conn.Exec(query, fullLink, shortLink)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}

	return models.Link{FullLink: fullLink, ShortLink: shortLink}, http.StatusCreated
}

func (dbs *PSQLStorage) GetURL(shortLink string) (interface{}, int) {
	conn, err := dbs.pool.Acquire()
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}
	defer dbs.pool.Release(conn)

	selectQuery := "select full_link from links where short_link = $1"
	row := conn.QueryRow(selectQuery, shortLink)
	var fullLink string
	err = row.Scan(&fullLink)
	if err != nil {
		return nil, http.StatusNotFound
	} else {
		return models.Link{FullLink: fullLink, ShortLink: shortLink}, http.StatusOK
	}
}
