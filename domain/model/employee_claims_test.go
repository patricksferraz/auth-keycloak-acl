package model_test

import (
	"testing"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_NewEmployee(t *testing.T) {

	id := uuid.NewV4().String()
	count := faker.Number().NumberInt(2)

	var roles []string
	for i := 0; i < count; i++ {
		roles = append(roles, faker.Lorem().Word())
	}

	employee, err := model.NewEmployeeClaims(id, roles)

	require.Nil(t, err)
	require.NotEmpty(t, uuid.FromStringOrNil(employee.ID))
	require.Equal(t, employee.Roles, roles)

	_, err = model.NewEmployeeClaims("", roles)
	require.NotNil(t, err)
}
