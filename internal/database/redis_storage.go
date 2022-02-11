package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ronmount/ozon_go/internal/config_parser"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/ronmount/ozon_go/internal/tools"
	"log"
	"net/http"
)

type RedisStorage struct {
	ctx        context.Context
	fullAsKey  *redis.Client
	shortAsKey *redis.Client
}

func NewRedisStorage() *RedisStorage {
	rs := &RedisStorage{}
	rs.ctx = context.Background()
	redisURL := config_parser.LoadRedisConfigs()
	rs.fullAsKey = redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})
	rs.shortAsKey = redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   1,
	})
	if e1, e2 := rs.fullAsKey.Ping(rs.ctx).Err(), rs.shortAsKey.Ping(rs.ctx).Err(); e1 != nil || e2 != nil {
		log.Fatalln("Connecting to Redis failed.")
	} else {
		log.Println("Redis successfully connected.")
	}

	return rs
}

func (rs *RedisStorage) AddURL(fullLink string) (interface{}, int) {
	var (
		response interface{}
		status   int
	)

	if alreadySavedLink, err := rs.fullAsKey.Get(rs.ctx, fullLink).Result(); err == nil {
		response, status = models.Link{FullLink: fullLink, ShortLink: alreadySavedLink}, http.StatusConflict
	} else if shortLink, err := tools.GenerateToken(); err == nil {
		rs.fullAsKey.Set(rs.ctx, fullLink, shortLink, 0)
		rs.shortAsKey.Set(rs.ctx, shortLink, fullLink, 0)
		response, status = models.Link{FullLink: fullLink, ShortLink: shortLink}, http.StatusCreated
	} else {
		response, status = nil, http.StatusInternalServerError
	}

	return response, status
}

func (rs *RedisStorage) GetURL(shortLink string) (interface{}, int) {
	var (
		response interface{}
		status   int
	)

	if fullLink, err := rs.shortAsKey.Get(rs.ctx, shortLink).Result(); err == nil {
		response, status = models.Link{FullLink: fullLink, ShortLink: shortLink}, http.StatusOK
	} else {
		response, status = nil, http.StatusNotFound
	}

	return response, status
}
