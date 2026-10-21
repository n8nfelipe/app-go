package usecase

import (
	"errors"

	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/repositorie"
	"github.com/n8nfelipe/app-go.git/internal/security"
)

var ErrInvalidCredentials = errors.New("credenciais inválidas")

type UserUseCase struct {
	UserRepository *repositorie.UserRepository
}

func NewUserUseCase(userRepository *repositorie.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (uc *UserUseCase) GetByUsernameAndPassword(username, password string) (*model.User, error) {
	user := uc.UserRepository.GetByUsername(username)
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	ok, err := security.VerifyPassword(password, user.Password)
	if err != nil || !ok {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (uc *UserUseCase) Create(user model.User) (uint, error) {
	hash, err := security.GeneratePassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hash
	return uc.UserRepository.Create(user)
}
