package model_test

import (
	"testing"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/model"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestModel_NewJWT(t *testing.T) {

	accessToken := faker.Lorem().Characters(17)
	idToken := faker.Lorem().Characters(17)
	expiresIn := faker.Number().NumberInt(3)
	refreshExpiresIn := faker.Number().NumberInt(3)
	refreshToken := faker.Lorem().Characters(17)
	tokenType := faker.Lorem().Word()
	notBeforePolicy := faker.Number().NumberInt(3)
	sessionState := faker.Lorem().Word()
	scope := faker.Lorem().Word()

	jwt, err := model.NewJWT(accessToken, idToken, refreshToken, tokenType, sessionState, scope, expiresIn, refreshExpiresIn, notBeforePolicy)

	require.Nil(t, err)
	require.Equal(t, jwt.AccessToken, accessToken)
	require.Equal(t, jwt.IDToken, idToken)
	require.Equal(t, jwt.RefreshToken, refreshToken)
	require.Equal(t, jwt.TokenType, tokenType)
	require.Equal(t, jwt.SessionState, sessionState)
	require.Equal(t, jwt.Scope, scope)
	require.Equal(t, jwt.ExpiresIn, expiresIn)
	require.Equal(t, jwt.RefreshExpiresIn, refreshExpiresIn)
	require.Equal(t, jwt.NotBeforePolicy, notBeforePolicy)

	_, err = model.NewJWT("", idToken, refreshToken, tokenType, sessionState, scope, expiresIn, refreshExpiresIn, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, "", refreshToken, tokenType, sessionState, scope, expiresIn, refreshExpiresIn, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, idToken, "", tokenType, sessionState, scope, expiresIn, refreshExpiresIn, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, idToken, refreshToken, "", sessionState, scope, expiresIn, refreshExpiresIn, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, idToken, refreshToken, tokenType, "", scope, expiresIn, refreshExpiresIn, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, idToken, refreshToken, tokenType, sessionState, "", expiresIn, refreshExpiresIn, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, idToken, refreshToken, tokenType, sessionState, scope, 0, refreshExpiresIn, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, idToken, refreshToken, tokenType, sessionState, scope, expiresIn, 0, notBeforePolicy)
	require.NotNil(t, err)
	_, err = model.NewJWT(accessToken, idToken, refreshToken, tokenType, sessionState, scope, expiresIn, refreshExpiresIn, 0)
	require.NotNil(t, err)
}
