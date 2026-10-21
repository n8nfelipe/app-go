package db

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type migration struct {
	Version int
	UpSQL   string
	DownSQL string
}

const createVersionsTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version    BIGINT PRIMARY KEY,
    dirty      BOOLEAN NOT NULL DEFAULT false
)`

// Up aplica todas as migrations pendentes, cada uma em sua própria transação.
func Up() error {
	conn, err := OpenConnection()
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if _, err := conn.Exec(createVersionsTable); err != nil {
		return fmt.Errorf("criando tabela de controle de migrations: %w", err)
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	applied, err := appliedVersions(conn)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}
		if err := applyMigration(conn, m); err != nil {
			return fmt.Errorf("migration %d: %w", m.Version, err)
		}
	}
	return nil
}

// Down reverte a última migration aplicada.
func Down() error {
	conn, err := OpenConnection()
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if _, err := conn.Exec(createVersionsTable); err != nil {
		return fmt.Errorf("criando tabela de controle de migrations: %w", err)
	}

	last, found, err := lastAppliedVersion(conn)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("nenhuma migration aplicada para reverter")
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.Version != last {
			continue
		}
		return revertMigration(conn, m)
	}
	return fmt.Errorf("migration aplicada %d não encontrada nos arquivos", last)
}

// Version retorna a última migration aplicada (0 se nenhuma).
func Version() (version int, dirty bool, err error) {
	conn, err := OpenConnection()
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = conn.Close() }()

	if _, err := conn.Exec(createVersionsTable); err != nil {
		return 0, false, fmt.Errorf("criando tabela de controle de migrations: %w", err)
	}

	last, found, err := lastAppliedVersion(conn)
	if err != nil {
		return 0, false, err
	}
	if !found {
		return 0, false, nil
	}
	// Cada migration roda em transação, então nunca fica "dirty".
	return last, false, nil
}

// loadMigrations lê os arquivos embedados no formato NNN_descricao.{up,down}.sql.
func loadMigrations() ([]migration, error) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("lendo diretório de migrations: %w", err)
	}

	byVersion := map[int]*migration{}
	for _, e := range entries {
		name := e.Name()
		switch {
		case strings.HasSuffix(name, ".up.sql"):
			m, err := parseMigrationName(name, ".up.sql")
			if err != nil {
				return nil, err
			}
			content, err := migrationsFS.ReadFile("migrations/" + name)
			if err != nil {
				return nil, fmt.Errorf("lendo %s: %w", name, err)
			}
			if byVersion[m.Version] == nil {
				byVersion[m.Version] = &migration{Version: m.Version}
			}
			byVersion[m.Version].UpSQL = string(content)
		case strings.HasSuffix(name, ".down.sql"):
			m, err := parseMigrationName(name, ".down.sql")
			if err != nil {
				return nil, err
			}
			content, err := migrationsFS.ReadFile("migrations/" + name)
			if err != nil {
				return nil, fmt.Errorf("lendo %s: %w", name, err)
			}
			if byVersion[m.Version] == nil {
				byVersion[m.Version] = &migration{Version: m.Version}
			}
			byVersion[m.Version].DownSQL = string(content)
		}
	}

	migrations := make([]migration, 0, len(byVersion))
	for _, m := range byVersion {
		if m.UpSQL == "" {
			return nil, fmt.Errorf("migration %d sem arquivo .up.sql", m.Version)
		}
		if m.DownSQL == "" {
			return nil, fmt.Errorf("migration %d sem arquivo .down.sql", m.Version)
		}
		migrations = append(migrations, *m)
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].Version < migrations[j].Version })

	return migrations, nil
}

func parseMigrationName(name, suffix string) (migration, error) {
	base := strings.TrimSuffix(name, suffix)
	num := strings.SplitN(base, "_", 2)[0]
	version, err := strconv.Atoi(num)
	if err != nil {
		return migration{}, fmt.Errorf("nome de migration inválido %q (esperado NNN_descricao%s): %w", name, suffix, err)
	}
	return migration{Version: version}, nil
}

func applyMigration(conn *sql.DB, m migration) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(m.UpSQL); err != nil {
		return fmt.Errorf("executando .up.sql: %w", err)
	}
	if _, err := tx.Exec("INSERT INTO schema_migrations (version, dirty) VALUES ($1, false)", m.Version); err != nil {
		return fmt.Errorf("registrando versão: %w", err)
	}
	return tx.Commit()
}

func revertMigration(conn *sql.DB, m migration) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(m.DownSQL); err != nil {
		return fmt.Errorf("executando .down.sql: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", m.Version); err != nil {
		return fmt.Errorf("removendo registro da versão: %w", err)
	}
	return tx.Commit()
}

func appliedVersions(conn *sql.DB) (map[int]bool, error) {
	rows, err := conn.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("consultando migrations aplicadas: %w", err)
	}
	defer func() { _ = rows.Close() }()

	applied := map[int]bool{}
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("lendo versão aplicada: %w", err)
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

func lastAppliedVersion(conn *sql.DB) (int, bool, error) {
	var version int
	err := conn.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("consultando última versão: %w", err)
	}
	return version, true, nil
}
