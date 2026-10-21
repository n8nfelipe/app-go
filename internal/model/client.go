package model

import (
	"database/sql"
	"log"

	"github.com/n8nfelipe/app-go.git/db"
)

type Client struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UserID    uint   `json:"user_id"`
}

func (c *Client) GetByID(id int, userID uint) (*Client, error) {
	conn, err := db.Pool()
	if err != nil {
		log.Printf("falha ao obter pool de conexões: %v", err)
		return nil, err
	}

	sqlQuery := "SELECT id, name, email, created_at, user_id FROM clients WHERE id = $1 AND user_id = $2"

	err = conn.QueryRow(sqlQuery, id, userID).Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt, &c.UserID)
	if err != nil {
		return nil, err // sql.ErrNoRows é propagado tal qual
	}
	return c, nil
}

func (c *Client) GetAll(userID uint) ([]Client, error) {
	conn, err := db.Pool()
	if err != nil {
		log.Printf("falha ao obter pool de conexões: %v", err)
		return nil, err
	}

	sqlQuery := "SELECT id, name, email, created_at, user_id FROM clients WHERE user_id = $1"

	rows, err := conn.Query(sqlQuery, userID)
	if err != nil {
		log.Printf("falha ao listar clients: %v", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt, &client.UserID); err != nil {
			log.Printf("falha ao ler client: %v", err)
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, rows.Err()
}

func (c *Client) Update(client *Client) error {
	conn, err := db.Pool()
	if err != nil {
		log.Printf("falha ao obter pool de conexões: %v", err)
		return err
	}

	sqlQuery := "UPDATE clients SET name = $1, email = $2 WHERE id = $3 AND user_id = $4"

	res, err := conn.Exec(sqlQuery, client.Name, client.Email, client.ID, client.UserID)
	if err != nil {
		log.Printf("falha ao atualizar client: %v", err)
		return err
	}
	return requireAffected(res)
}

func (c *Client) Create(client *Client) (id int, err error) {
	conn, err := db.Pool()
	if err != nil {
		log.Printf("falha ao obter pool de conexões: %v", err)
		return 0, err
	}

	sqlQuery := "INSERT INTO clients (name, email, user_id) VALUES ($1, $2, $3) RETURNING id"

	err = conn.QueryRow(sqlQuery, client.Name, client.Email, client.UserID).Scan(&client.ID)
	if err != nil {
		log.Printf("falha ao criar client: %v", err)
		return 0, err
	}
	return client.ID, nil
}

func (c *Client) Delete(id int, userID uint) error {
	conn, err := db.Pool()
	if err != nil {
		log.Printf("falha ao obter pool de conexões: %v", err)
		return err
	}

	res, err := conn.Exec("DELETE FROM clients WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		log.Printf("falha ao excluir client: %v", err)
		return err
	}
	return requireAffected(res)
}

func requireAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
