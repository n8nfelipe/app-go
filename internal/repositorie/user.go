package repositorie

import "github.com/n8nfelipe/app-go.git/internal/model"

type UserRepository struct {
	User *model.User
}

func NewUserRepository(user *model.User) *UserRepository {
	return &UserRepository{
		User: user,
	}
}

func (r *UserRepository) GetByUsername(username string) *model.User {
	return r.User.GetByUsername(username)
}

func (r *UserRepository) Create(user model.User) (uint, error) {
	return r.User.Create(user)
}
