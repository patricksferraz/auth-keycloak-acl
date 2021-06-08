package utils_test

import (
	"os"
	"testing"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/utils"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestUtils_GetEnv(t *testing.T) {

	key := faker.Lorem().Word()
	defaultVal := faker.Lorem().Word()

	result := utils.GetEnv(key, defaultVal)
	require.Equal(t, result, defaultVal)

	otherVal := faker.Lorem().Word()
	os.Setenv(key, otherVal)

	result = utils.GetEnv(key, defaultVal)
	require.Equal(t, result, otherVal)
}
