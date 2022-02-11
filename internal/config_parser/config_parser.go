package config_parser

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func LoadRedisConfigs() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	RedisHost := os.Getenv("REDIS_HOST")
	RedisPort := os.Getenv("REDIS_PORT")
	if len(RedisHost) == 0 || len(RedisPort) == 0 || !isNumeric(RedisPort) {
		log.Fatal("Broken env.")
	}
	config := fmt.Sprintf("%s:%s", RedisHost, RedisPort)
	return config
}

func LoadPSQLConfigs() *pgx.ConnPoolConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	PgHost := os.Getenv("POSTGRES_HOST")
	PgPort := os.Getenv("POSTGRES_PORT")
	PgDatabase := os.Getenv("POSTGRES_DB")
	PgUser := os.Getenv("POSTGRES_USER")
	PgPassword := os.Getenv("POSTGRES_PASSWORD")
	if len(PgHost) == 0 || len(PgPort) == 0 || len(PgDatabase) == 0 ||
		len(PgUser) == 0 || len(PgPassword) == 0 || !isNumeric(PgPort) {
		log.Fatal("Broken env.")
	}
	PgPortInt, _ := strconv.ParseUint(PgPort, 10, 64)
	config := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     PgHost,
			Port:     uint16(PgPortInt),
			Database: PgDatabase,
			User:     PgUser,
			Password: PgPassword,
		}}
	return &config
}
