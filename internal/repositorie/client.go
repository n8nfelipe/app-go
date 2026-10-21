package repositorie

import "github.com/n8nfelipe/app-go.git/internal/model"

type ClientRepository struct {
	Client *model.Client
}

func NewClientRepository(client *model.Client) *ClientRepository {
	return &ClientRepository{
		Client: client,
	}
}

func (r *ClientRepository) GetByID(id int) *model.Client {
	return r.Client.GetByID(id)
}

func (r *ClientRepository) GetAll() []model.Client {
	clients, err := r.Client.GetAll()
	if err != nil {
		return nil
	}
	return clients
}

func (r *ClientRepository) Update(client *model.Client) error {
	return r.Client.Update(client)
}

func (r *ClientRepository) Create(client *model.Client) (id int, err error) {
	return r.Client.Create(client)
}
