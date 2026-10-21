package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dbEnvKeys = []string{
	"APP_DB_HOST",
	"APP_DB_PORT",
	"APP_DB_USER",
	"APP_DB_PASSWORD",
	"APP_DB_DBNAME",
}

func TestLoad_LeVariaveisDeAmbiente(t *testing.T) {
	t.Setenv("APP_DB_HOST", "db-exemplo")
	t.Setenv("APP_DB_PORT", "6543")
	t.Setenv("APP_DB_USER", "app")
	t.Setenv("APP_DB_PASSWORD", "segredo")
	t.Setenv("APP_DB_DBNAME", "clientes")
	t.Setenv("APP_SERVER_PORT", "9999")

	cfg, err := Load()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "db-exemplo", cfg.DB.Host)
	assert.Equal(t, "6543", cfg.DB.Port)
	assert.Equal(t, "app", cfg.DB.User)
	assert.Equal(t, "segredo", cfg.DB.Password)
	assert.Equal(t, "clientes", cfg.DB.DBName)
	assert.Equal(t, "9999", cfg.Server.Port)
}

func TestLoad_UseDefaultsQuandoEnvAusente(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, key := range dbEnvKeys {
		t.Setenv(key, "")
	}
	t.Setenv("APP_SERVER_PORT", "")

	cfg, err := Load()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, "postgres", cfg.DB.User)
	assert.Equal(t, "password", cfg.DB.Password)
	assert.Equal(t, "mydb", cfg.DB.DBName)
	assert.Equal(t, "8081", cfg.Server.Port)
}
