package usecase

import (
	"errors"
	"strings"
	"testing"

	"github.com/n8nfelipe/app-go.git/internal/apperror"
	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/n8nfelipe/app-go.git/internal/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepository struct {
	getByUsername func(username string) (*model.User, error)
	create        func(user model.User) (uint, error)

	getByUsernameCalls []string
	created            []model.User
}

func (m *mockUserRepository) GetByUsername(username string) (*model.User, error) {
	m.getByUsernameCalls = append(m.getByUsernameCalls, username)
	return m.getByUsername(username)
}

func (m *mockUserRepository) Create(user model.User) (uint, error) {
	m.created = append(m.created, user)
	return m.create(user)
}

func TestUserUseCase_GetByUsernameAndPassword_UsuarioNaoExiste(t *testing.T) {
	repo := &mockUserRepository{getByUsername: func(username string) (*model.User, error) { return nil, nil }}
	uc := NewUserUseCase(repo)

	user, err := uc.GetByUsernameAndPassword("fulano", "senha")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Equal(t, []string{"fulano"}, repo.getByUsernameCalls)
}

func TestUserUseCase_GetByUsernameAndPassword_ErroDeRepositorio(t *testing.T) {
	repo := &mockUserRepository{
		getByUsername: func(username string) (*model.User, error) {
			return nil, errors.New("banco indisponível")
		},
	}
	uc := NewUserUseCase(repo)

	user, err := uc.GetByUsernameAndPassword("felipe", "senha")

	assert.Nil(t, user)
	require.Error(t, err)
	var applicationError apperror.Error
	require.ErrorAs(t, err, &applicationError)
	assert.Equal(t, apperror.CodeInternal, applicationError.Code())
	assert.NotErrorIs(t, err, ErrInvalidCredentials)
}

func TestUserUseCase_GetByUsernameAndPassword_SenhaErrada(t *testing.T) {
	hash, err := security.GeneratePassword("senha-correta")
	require.NoError(t, err)

	repo := &mockUserRepository{
		getByUsername: func(username string) (*model.User, error) {
			return &model.User{ID: 1, Username: "felipe", Password: hash}, nil
		},
	}
	uc := NewUserUseCase(repo)

	user, err := uc.GetByUsernameAndPassword("felipe", "senha-errada")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestUserUseCase_GetByUsernameAndPassword_HashCorrompido(t *testing.T) {
	repo := &mockUserRepository{
		getByUsername: func(username string) (*model.User, error) {
			return &model.User{ID: 1, Username: "felipe", Password: "nao-e-hash"}, nil
		},
	}
	uc := NewUserUseCase(repo)

	user, err := uc.GetByUsernameAndPassword("felipe", "qualquer")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestUserUseCase_GetByUsernameAndPassword_Sucesso(t *testing.T) {
	hash, err := security.GeneratePassword("senha-correta")
	require.NoError(t, err)

	repo := &mockUserRepository{
		getByUsername: func(username string) (*model.User, error) {
			return &model.User{ID: 42, Username: "felipe", Password: hash}, nil
		},
	}
	uc := NewUserUseCase(repo)

	user, err := uc.GetByUsernameAndPassword("felipe", "senha-correta")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, uint(42), user.ID)
	assert.Equal(t, "felipe", user.Username)
}

func TestUserUseCase_Create_GuardaSenhaHasheada(t *testing.T) {
	repo := &mockUserRepository{
		create: func(user model.User) (uint, error) { return 5, nil },
	}
	uc := NewUserUseCase(repo)

	id, err := uc.Create(model.User{Username: "novo", Password: "minha-senha-plaintext"})

	require.NoError(t, err)
	assert.Equal(t, uint(5), id)
	require.Len(t, repo.created, 1)

	stored := repo.created[0]
	assert.NotEqual(t, "minha-senha-plaintext", stored.Password)
	assert.True(t, strings.HasPrefix(stored.Password, "$argon2id$"))

	ok, err := security.VerifyPassword("minha-senha-plaintext", stored.Password)
	require.NoError(t, err)
	assert.True(t, ok)
}
