package usecase

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/n8nfelipe/app-go.git/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClientRepository struct {
	getByID func(id int, userID uint) (*model.Client, error)
	getAll  func(userID uint) ([]model.Client, error)
	update  func(client *model.Client) error
	create  func(client *model.Client) (int, error)
	del     func(id int, userID uint) error

	getByIDCalls [][2]int
	getAllCalls  []uint
	updateCalls  []*model.Client
	createCalls  []*model.Client
	deleteCalls  [][2]int
}

func (m *mockClientRepository) GetByID(id int, userID uint) (*model.Client, error) {
	m.getByIDCalls = append(m.getByIDCalls, [2]int{id, int(userID)})
	return m.getByID(id, userID)
}

func (m *mockClientRepository) GetAll(userID uint) ([]model.Client, error) {
	m.getAllCalls = append(m.getAllCalls, userID)
	return m.getAll(userID)
}

func (m *mockClientRepository) Update(client *model.Client) error {
	m.updateCalls = append(m.updateCalls, client)
	return m.update(client)
}

func (m *mockClientRepository) Create(client *model.Client) (int, error) {
	m.createCalls = append(m.createCalls, client)
	return m.create(client)
}

func (m *mockClientRepository) Delete(id int, userID uint) error {
	m.deleteCalls = append(m.deleteCalls, [2]int{id, int(userID)})
	return m.del(id, userID)
}

func TestClientUseCase_GetClientByID_PassaEscopoDoUsuario(t *testing.T) {
	repo := &mockClientRepository{
		getByID: func(id int, userID uint) (*model.Client, error) {
			return &model.Client{ID: id, Name: "Empresa X", UserID: userID}, nil
		},
	}
	uc := NewClientUseCase(repo)

	client, err := uc.GetClientByID(7, 42)

	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, 7, client.ID)
	assert.Equal(t, uint(42), client.UserID)
	assert.Equal(t, [][2]int{{7, 42}}, repo.getByIDCalls)
}

func TestClientUseCase_GetClientByID_NaoEncontrado(t *testing.T) {
	repo := &mockClientRepository{
		getByID: func(id int, userID uint) (*model.Client, error) {
			return nil, sql.ErrNoRows
		},
	}
	uc := NewClientUseCase(repo)

	client, err := uc.GetClientByID(999, 42)

	assert.Nil(t, client)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestClientUseCase_GetClientByID_ErroDeRepositorio(t *testing.T) {
	repo := &mockClientRepository{
		getByID: func(id int, userID uint) (*model.Client, error) {
			return nil, errors.New("banco indisponível")
		},
	}
	uc := NewClientUseCase(repo)

	_, err := uc.GetClientByID(1, 42)

	require.Error(t, err)
	assert.NotErrorIs(t, err, sql.ErrNoRows)
}

func TestClientUseCase_GetClients_PassaEscopoDoUsuario(t *testing.T) {
	want := []model.Client{
		{ID: 1, Name: "A", UserID: 42},
		{ID: 2, Name: "B", UserID: 42},
	}
	repo := &mockClientRepository{getAll: func(userID uint) ([]model.Client, error) { return want, nil }}
	uc := NewClientUseCase(repo)

	got, err := uc.GetClients(42)

	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, []uint{42}, repo.getAllCalls)
}

func TestClientUseCase_GetClients_Erro(t *testing.T) {
	repo := &mockClientRepository{
		getAll: func(userID uint) ([]model.Client, error) { return nil, errors.New("timeout") },
	}
	uc := NewClientUseCase(repo)

	clients, err := uc.GetClients(42)

	assert.Nil(t, clients)
	require.Error(t, err)
}

func TestClientUseCase_Create_PreservaDono(t *testing.T) {
	repo := &mockClientRepository{
		create: func(client *model.Client) (int, error) { return 10, nil },
	}
	uc := NewClientUseCase(repo)

	id, err := uc.CreateClient(&model.Client{Name: "Nova", Email: "nova@empresa.com", UserID: 42})

	require.NoError(t, err)
	assert.Equal(t, 10, id)
	require.Len(t, repo.createCalls, 1)
	assert.Equal(t, uint(42), repo.createCalls[0].UserID)
}

func TestClientUseCase_Create_Erro(t *testing.T) {
	repo := &mockClientRepository{
		create: func(client *model.Client) (int, error) { return 0, errors.New("banco indisponível") },
	}
	uc := NewClientUseCase(repo)

	id, err := uc.CreateClient(&model.Client{Name: "Nova"})

	require.Error(t, err)
	assert.Equal(t, 0, id)
}

func TestClientUseCase_Update_Successo(t *testing.T) {
	repo := &mockClientRepository{update: func(client *model.Client) error { return nil }}
	uc := NewClientUseCase(repo)

	err := uc.UpdateClient(&model.Client{ID: 1, Name: "Atualizado", UserID: 42})

	require.NoError(t, err)
	require.Len(t, repo.updateCalls, 1)
	assert.Equal(t, "Atualizado", repo.updateCalls[0].Name)
}

func TestClientUseCase_Update_NaoEncontrado(t *testing.T) {
	repo := &mockClientRepository{update: func(client *model.Client) error { return sql.ErrNoRows }}
	uc := NewClientUseCase(repo)

	err := uc.UpdateClient(&model.Client{ID: 404})

	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestClientUseCase_Delete_PassaEscopoDoUsuario(t *testing.T) {
	repo := &mockClientRepository{del: func(id int, userID uint) error { return nil }}
	uc := NewClientUseCase(repo)

	err := uc.DeleteClient(3, 42)

	require.NoError(t, err)
	assert.Equal(t, [][2]int{{3, 42}}, repo.deleteCalls)
}

func TestClientUseCase_Delete_NaoEncontrado(t *testing.T) {
	repo := &mockClientRepository{del: func(id int, userID uint) error { return sql.ErrNoRows }}
	uc := NewClientUseCase(repo)

	err := uc.DeleteClient(3, 42)

	assert.ErrorIs(t, err, sql.ErrNoRows)
}
