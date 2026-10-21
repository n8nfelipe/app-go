package model

import (
	"log"

	"github.com/n8nfelipe/app-go.git/db"
)

type Client struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (c *Client) GetByID(id int) *Client {
	db, err := db.OpenConnection()
	if err != nil {
		// Lidar com o erro de conexão ao banco de dados
		log.Printf("Failed to update client: %v", err)
		return nil
	}
	defer db.Close()

	sql := "SELECT id, name, email, created_at FROM clients WHERE id = $1"

	row := db.QueryRow(sql, id)

	err = row.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
	if err != nil {
		// Lidar com o erro de consulta ao banco de dados
		log.Printf("Failed to update client: %v", err)
		return nil
	}
	return c
}

func (c *Client) GetAll() (clients []Client, err error) {
	db, err := db.OpenConnection()
	if err != nil {
		// Lidar com o erro de conexão ao banco de dados
		return
	}
	defer db.Close()

	sql := "SELECT id, name, email, created_at FROM clients"

	row, err := db.Query(sql)
	if err != nil {
		// Lidar com o erro de consulta ao banco de dados
		log.Printf("Failed to update client: %v", err)
		return clients, err
	}

	for row.Next() {
		var client Client
		err = row.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt)
		if err != nil {
			// Lidar com o erro de consulta ao banco de dados
			continue
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func (c *Client) Update(client *Client) error {
	db, err := db.OpenConnection()
	if err != nil {
		// Lidar com o erro de conexão ao banco de dados
		return err
	}
	defer db.Close()

	sql := "UPDATE clients SET name = $1, email = $2 WHERE id = $3"

	_, err = db.Exec(sql, client.Name, client.Email, client.ID)
	if err != nil {
		// Lidar com o erro de atualização no banco de dados
		log.Printf("Failed to update client: %v", err)
		return err
	}
	return nil
}

func (c *Client) Create(client *Client) (id int, err error) {
	db, err := db.OpenConnection()
	if err != nil {
		log.Printf("Failed to update client: %v", err)
		// Lidar com o erro de conexão ao banco de dados
		return 0, err
	}
	defer db.Close()

	sql := "INSERT INTO clients (name, email) VALUES ($1, $2) RETURNING id"
	err = db.QueryRow(sql, client.Name, client.Email).Scan(&client.ID)
	if err != nil {
		// Lidar com o erro de inserção no banco de dados
		log.Printf("Failed to update client: %v", err)
		return 0, err
	}
	return client.ID, nil
}
