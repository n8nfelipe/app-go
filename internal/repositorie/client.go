package repositorie

import (
	"database/sql"
	"errors"

	"github.com/n8nfelipe/app-go.git/internal/apperror"
	"github.com/n8nfelipe/app-go.git/internal/model"
)

type ClientRepository struct {
	Client *model.Client
}

func NewClientRepository(client *model.Client) *ClientRepository {
	return &ClientRepository{
		Client: client,
	}
}

func (r *ClientRepository) GetByID(id int, userID uint) (*model.Client, error) {
	client, err := r.Client.GetByID(id, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.New(apperror.CodeNotFound, "Cliente não encontrado", err)
	}
	if err != nil {
		return nil, apperror.Wrap(err, "buscando cliente")
	}
	return client, nil
}

func (r *ClientRepository) GetAll(userID uint) ([]model.Client, error) {
	clients, err := r.Client.GetAll(userID)
	if err != nil {
		return nil, apperror.Wrap(err, "listando clientes")
	}
	return clients, nil
}

func (r *ClientRepository) Update(client *model.Client) error {
	err := r.Client.Update(client)
	if errors.Is(err, sql.ErrNoRows) {
		return apperror.New(apperror.CodeNotFound, "Cliente não encontrado", err)
	}
	return apperror.Wrap(err, "atualizando cliente")
}

func (r *ClientRepository) Create(client *model.Client) (id int, err error) {
	id, err = r.Client.Create(client)
	if err != nil {
		return 0, apperror.Wrap(err, "criando cliente")
	}
	return id, nil
}

func (r *ClientRepository) Delete(id int, userID uint) error {
	err := r.Client.Delete(id, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return apperror.New(apperror.CodeNotFound, "Cliente não encontrado", err)
	}
	return apperror.Wrap(err, "excluindo cliente")
}
