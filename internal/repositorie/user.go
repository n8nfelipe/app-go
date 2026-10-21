package repositorie

import (
	"database/sql"
	"errors"

	"github.com/n8nfelipe/app-go.git/internal/apperror"
	"github.com/n8nfelipe/app-go.git/internal/model"
)

type UserRepository struct {
	User *model.User
}

func NewUserRepository(user *model.User) *UserRepository {
	return &UserRepository{
		User: user,
	}
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	user, err := r.User.GetByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.New(apperror.CodeNotFound, "Usuário não encontrado", err)
	}
	if err != nil {
		return nil, apperror.Wrap(err, "buscando usuário")
	}
	return user, nil
}

func (r *UserRepository) Create(user model.User) (uint, error) {
	id, err := r.User.Create(user)
	if err != nil {
		return 0, apperror.Wrap(err, "criando usuário")
	}
	return id, nil
}
