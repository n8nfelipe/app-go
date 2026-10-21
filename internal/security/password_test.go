package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePassword_GeraHashArgon2id(t *testing.T) {
	hash, err := GeneratePassword("minha-senha")

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(hash, "$argon2id$v="))
}

func TestGeneratePassword_SaltUnicoPorChamada(t *testing.T) {
	hash1, err := GeneratePassword("mesma-senha")
	require.NoError(t, err)
	hash2, err := GeneratePassword("mesma-senha")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2)
}

func TestVerifyPassword_SenhaCorreta(t *testing.T) {
	hash, err := GeneratePassword("s3nh4")
	require.NoError(t, err)

	ok, err := VerifyPassword("s3nh4", hash)

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestVerifyPassword_SenhaErrada(t *testing.T) {
	hash, err := GeneratePassword("s3nh4")
	require.NoError(t, err)

	ok, err := VerifyPassword("outra-senha", hash)

	require.NoError(t, err)
	assert.False(t, ok)
}

func TestVerifyPassword_HashInvalido(t *testing.T) {
	ok, err := VerifyPassword("qualquer", "nao-e-um-hash")

	assert.False(t, ok)
	assert.ErrorIs(t, err, ErrInvalidHash)
}

func TestVerifyPassword_AlgoritmoDiferente(t *testing.T) {
	ok, err := VerifyPassword("qualquer", "$argon2i$v=19$m=65536,t=3,p=4$c2FsdHNhbHQ$aGFzaGhhc2g")

	assert.False(t, ok)
	assert.ErrorIs(t, err, ErrInvalidHash)
}

func TestVerifyPassword_VersaoIncompativel(t *testing.T) {
	ok, err := VerifyPassword("qualquer", "$argon2id$v=1$m=65536,t=3,p=4$c2FsdHNhbHQ$aGFzaGhhc2g")

	assert.False(t, ok)
	assert.ErrorIs(t, err, ErrIncompatibleVersion)
}
