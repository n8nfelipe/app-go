package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Importa o driver do PostgreSQL
	"github.com/n8nfelipe/app-go.git/config"
)

func OpenConnection() (*sql.DB, error) {
	cfg := config.GetDbConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco de dados: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão com o banco de dados: %w", err)
	}

	return db, nil
}
