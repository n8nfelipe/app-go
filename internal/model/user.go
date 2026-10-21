package model

import (
	"log"

	"github.com/n8nfelipe/app-go.git/db"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *User) Create(user User) (uint, error) {
	conn, err := db.Pool()
	if err != nil {
		log.Printf("falha ao obter pool de conexões: %v", err)
		return 0, err
	}

	sql := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	var id uint
	err = conn.QueryRow(sql, user.Username, user.Password).Scan(&id)
	if err != nil {
		log.Printf("falha ao inserir usuário: %v", err)
		return 0, err
	}
	return id, nil
}

func (c *User) GetByUsername(username string) (*User, error) {
	conn, err := db.Pool()
	if err != nil {
		log.Printf("falha ao obter pool de conexões: %v", err)
		return nil, err
	}

	sql := "SELECT id, username, password FROM users WHERE username = $1"

	err = conn.QueryRow(sql, username).Scan(&c.ID, &c.Username, &c.Password)
	if err != nil {
		return nil, err
	}
	return c, nil
}
