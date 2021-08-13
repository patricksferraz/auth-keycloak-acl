package entity_test

import (
	"testing"

	"github.com/c-4u/auth-service/domain/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_NewClaims(t *testing.T) {

	userID := uuid.NewV4().String()
	employeeID := uuid.NewV4().String()
	username := faker.Internet().UserName()
	count := faker.Number().NumberInt(2)

	var roles []string
	for i := 0; i < count; i++ {
		roles = append(roles, faker.Lorem().Word())
	}

	claims, err := entity.NewClaims(userID, username, employeeID, roles)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(claims.UserID))
	require.Equal(t, claims.Roles, roles)

	_, err = entity.NewClaims("", username, employeeID, roles)
	require.NotNil(t, err)
	_, err = entity.NewClaims(userID, "", employeeID, roles)
	require.NotNil(t, err)
	_, err = entity.NewClaims(userID, username, "", roles)
	require.NotNil(t, err)
}
