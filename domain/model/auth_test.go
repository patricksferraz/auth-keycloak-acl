package model_test

import (
	"testing"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/model"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_NewAuth(t *testing.T) {

	username := faker.Internet().UserName()
	password := faker.Internet().Password(8, 20)

	auth, err := model.NewAuth(username, password)

	require.Nil(t, err)
	require.Equal(t, auth.Username, username)
	require.Equal(t, auth.Password, password)

	_, err = model.NewAuth("", password)
	require.NotNil(t, err)
	_, err = model.NewAuth(username, "")
	require.NotNil(t, err)
}
