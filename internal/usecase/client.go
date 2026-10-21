package usecase

import (
	"github.com/n8nfelipe/app-go.git/internal/apperror"
	"github.com/n8nfelipe/app-go.git/internal/model"
)

type ClientRepository interface {
	GetByID(id int, userID uint) (*model.Client, error)
	GetAll(userID uint) ([]model.Client, error)
	Update(client *model.Client) error
	Create(client *model.Client) (id int, err error)
	Delete(id int, userID uint) error
}

type ClientUseCase struct {
	ClientRepository ClientRepository
}

func NewClientUseCase(clientRepository ClientRepository) *ClientUseCase {
	return &ClientUseCase{
		ClientRepository: clientRepository,
	}
}

func (uc *ClientUseCase) GetClientByID(id int, userID uint) (*model.Client, error) {
	client, err := uc.ClientRepository.GetByID(id, userID)
	if err != nil {
		return nil, apperror.Wrap(err, "buscando cliente")
	}
	return client, nil
}

func (uc *ClientUseCase) GetClients(userID uint) ([]model.Client, error) {
	clients, err := uc.ClientRepository.GetAll(userID)
	if err != nil {
		return nil, apperror.Wrap(err, "listando clientes")
	}
	return clients, nil
}

func (uc *ClientUseCase) UpdateClient(client *model.Client) error {
	return apperror.Wrap(uc.ClientRepository.Update(client), "atualizando cliente")
}

func (uc *ClientUseCase) CreateClient(client *model.Client) (id int, err error) {
	id, err = uc.ClientRepository.Create(client)
	if err != nil {
		return 0, apperror.Wrap(err, "criando cliente")
	}
	return id, nil
}

func (uc *ClientUseCase) DeleteClient(id int, userID uint) error {
	return apperror.Wrap(uc.ClientRepository.Delete(id, userID), "excluindo cliente")
}
