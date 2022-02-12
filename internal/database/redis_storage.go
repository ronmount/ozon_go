package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ronmount/ozon_go/internal/config_parser"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/ronmount/ozon_go/internal/my_errors"
	"github.com/ronmount/ozon_go/internal/tools"
)

type RedisStorage struct {
	ctx        context.Context
	fullAsKey  *redis.Client
	shortAsKey *redis.Client
}

func NewRedisStorage() (*RedisStorage, error) {
	rs := &RedisStorage{}
	rs.ctx = context.Background()
	redisURL, err := config_parser.LoadRedisConfigs()
	if err != nil {
		return nil, err
	}
	rs.fullAsKey = redis.NewClient(&redis.Options{
		Addr: redisURL.(string),
		DB:   0,
	})
	rs.shortAsKey = redis.NewClient(&redis.Options{
		Addr: redisURL.(string),
		DB:   1,
	})

	if err := rs.fullAsKey.Ping(rs.ctx).Err(); err != nil {
		return nil, err
	} else if err = rs.fullAsKey.Ping(rs.ctx).Err(); err != nil {
		return nil, err
	}
	return rs, nil
}

func (rs *RedisStorage) AddURL(fullLink string) (interface{}, error) {
	var (
		response interface{}
		err      error
	)

	if alreadySavedLink, e := rs.fullAsKey.Get(rs.ctx, fullLink).Result(); e == nil {
		response, err = models.Link{FullLink: fullLink, ShortLink: alreadySavedLink}, nil
	} else if shortLink, e := tools.GenerateToken(); e == nil {
		rs.fullAsKey.Set(rs.ctx, fullLink, shortLink, 0)
		rs.shortAsKey.Set(rs.ctx, shortLink, fullLink, 0)
		response, err = models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	} else {
		response, err = nil, my_errors.HTTP500{}
	}

	return response, err
}

func (rs *RedisStorage) GetURL(shortLink string) (interface{}, error) {
	var (
		response interface{}
		err      error
	)

	if fullLink, e := rs.shortAsKey.Get(rs.ctx, shortLink).Result(); e == nil {
		response, err = models.Link{FullLink: fullLink, ShortLink: shortLink}, nil
	} else {
		response, err = nil, my_errors.HTTP404{}
	}

	return response, err
}
