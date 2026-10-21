package usecase

import (
	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/repositorie"
)

type ClientUseCase struct {
	ClientRepository *repositorie.ClientRepository
}

func NewClientUseCase(clientRepository *repositorie.ClientRepository) *ClientUseCase {
	return &ClientUseCase{
		ClientRepository: clientRepository,
	}
}

func (uc *ClientUseCase) GetClientByID(id int) *model.Client {
	return uc.ClientRepository.GetByID(id)
}

func (uc *ClientUseCase) GetClients() []model.Client {
	return uc.ClientRepository.GetAll()
}

func (uc *ClientUseCase) UpdateClient(client *model.Client) error {
	return uc.ClientRepository.Update(client)
}

func (uc *ClientUseCase) CreateClient(client *model.Client) (id int, err error) {
	return uc.ClientRepository.Create(client)
}
