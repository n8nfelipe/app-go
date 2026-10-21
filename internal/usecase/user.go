package usecase

import (
	"errors"

	"github.com/n8nfelipe/app-go.git/internal/apperror"
	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/security"
)

var ErrInvalidCredentials = errors.New("credenciais inválidas")

type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
	Create(user model.User) (uint, error)
}

type UserUseCase struct {
	UserRepository UserRepository
}

func NewUserUseCase(userRepository UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (uc *UserUseCase) GetByUsernameAndPassword(username, password string) (*model.User, error) {
	user, err := uc.UserRepository.GetByUsername(username)
	if err != nil {
		var repositoryError apperror.Error
		if errors.As(err, &repositoryError) && repositoryError.Code() == apperror.CodeNotFound {
			return nil, invalidCredentialsError()
		}
		return nil, apperror.Wrap(err, "buscando usuário para autenticação")
	}
	if user == nil {
		return nil, invalidCredentialsError()
	}

	ok, err := security.VerifyPassword(password, user.Password)
	if err != nil || !ok {
		return nil, invalidCredentialsError()
	}

	return user, nil
}

func (uc *UserUseCase) Create(user model.User) (uint, error) {
	hash, err := security.GeneratePassword(user.Password)
	if err != nil {
		return 0, apperror.Wrap(err, "gerando hash da senha")
	}
	user.Password = hash
	id, err := uc.UserRepository.Create(user)
	if err != nil {
		return 0, apperror.Wrap(err, "criando usuário")
	}
	return id, nil
}

func invalidCredentialsError() error {
	return apperror.New(apperror.CodeUnauthorized, "Credenciais inválidas", ErrInvalidCredentials)
}
