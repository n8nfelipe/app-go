package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

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

var (
	poolOnce sync.Once
	poolDB   *sql.DB
	poolErr  error
)

// Pool retorna a conexão compartilhada da aplicação, com pool de conexões.
// Diferente de OpenConnection, não deve ser fechada pelos chamadores.
func Pool() (*sql.DB, error) {
	poolOnce.Do(func() {
		conn, err := OpenConnection()
		if err != nil {
			poolErr = err
			return
		}
		conn.SetMaxOpenConns(20)
		conn.SetMaxIdleConns(5)
		conn.SetConnMaxLifetime(30 * time.Minute)
		poolDB = conn
	})
	return poolDB, poolErr
}
