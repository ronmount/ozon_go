package config_parser

import (
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadPSQLConfigs(t *testing.T) {
	t.Log("When normal config")
	data := []byte("POSTGRES_HOST=localhost\nPOSTGRES_PORT=5432\nPOSTGRES_DB=shortener\nPOSTGRES_USER=user\nPOSTGRES_PASSWORD=password")
	os.WriteFile(".env", data, 0644)
	host, port, db, user, password := "localhost", uint16(5432), "shortener", "user", "password"
	expectedConfig := &pgx.ConnPoolConfig{ConnConfig: pgx.ConnConfig{
		Host: host, Port: port, Database: db, User: user, Password: password,
	}}
	actualConfig, err := LoadPSQLConfigs()
	require.NoError(t, err, "Should be no error")
	require.Equal(t, expectedConfig, actualConfig, "Should be equal to expectedConfig")
	os.Remove(".env")

	t.Log("When broken config")
	data = []byte("I am broken")
	os.WriteFile(".env", data, 0644)
	_, err = LoadPSQLConfigs()
	require.Error(t, err, "Should be an error")
	os.Remove(".env")

	t.Log("When no .env")
	_, err = LoadPSQLConfigs()
	require.Error(t, err, "Should be an error")
}

func TestLoadRedisConfigs(t *testing.T) {
	data := []byte("REDIS_HOST=localhost\nREDIS_PORT=6379")
	os.WriteFile(".env", data, 0644)
	expectedConfig := "localhost:6379"
	actualConfig, err := LoadRedisConfigs()
	require.NoError(t, err, "Should be no error")
	require.Equal(t, expectedConfig, actualConfig, "Should be equal to expectedConfig")
	os.Remove(".env")

	t.Log("When broken config")
	data = []byte("I am broken")
	os.WriteFile(".env", data, 0644)
	_, err = LoadPSQLConfigs()
	require.Error(t, err, "Should be an error")
	os.Remove(".env")

	t.Log("When no .env")
	_, err = LoadPSQLConfigs()
	require.Error(t, err, "Should be an error")
}
