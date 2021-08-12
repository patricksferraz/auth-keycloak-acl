package entity_test

import (
	"testing"

	"github.com/c-4u/auth-service/domain/entity"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_NewAuth(t *testing.T) {

	username := faker.Internet().UserName()
	password := faker.Internet().Password(8, 20)

	auth, err := entity.NewAuth(username, password)

	require.Nil(t, err)
	require.Equal(t, auth.Username, username)
	require.Equal(t, auth.Password, password)

	_, err = entity.NewAuth("", password)
	require.NotNil(t, err)
	_, err = entity.NewAuth(username, "")
	require.NotNil(t, err)
}
