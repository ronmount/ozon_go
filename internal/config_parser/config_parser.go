package config_parser

import (
	"errors"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func LoadPSQLConfigs() (*pgx.ConnPoolConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	PgHost := os.Getenv("POSTGRES_HOST")
	PgPort := os.Getenv("POSTGRES_PORT")
	PgDatabase := os.Getenv("POSTGRES_DB")
	PgUser := os.Getenv("POSTGRES_USER")
	PgPassword := os.Getenv("POSTGRES_PASSWORD")
	if len(PgHost) == 0 || len(PgPort) == 0 || len(PgDatabase) == 0 ||
		len(PgUser) == 0 || len(PgPassword) == 0 || !isNumeric(PgPort) {
		return nil, errors.New("broken env")
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
	return &config, nil
}
