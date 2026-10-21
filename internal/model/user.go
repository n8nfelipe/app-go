package model

import (
	"log"

	"github.com/n8nfelipe/app-go.git/db"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

func (c *User) Create(user User) (uint, error) {
	db, err := db.OpenConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return 0, err
	}
	defer db.Close()

	sql := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	var id uint
	err = db.QueryRow(sql, user.Username, user.Password).Scan(&id)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return 0, err
	}
	return id, nil
}

func (c *User) GetByUsername(username string) *User {
	db, err := db.OpenConnection()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil
	}
	defer db.Close()

	sql := "SELECT id, username, password FROM users WHERE username = $1"

	row := db.QueryRow(sql, username)

	err = row.Scan(&c.ID, &c.Username, &c.Password)
	if err != nil {
		log.Printf("Failed to query user: %v", err)
		return nil
	}
	return c
}
