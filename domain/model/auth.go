package model

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Auth struct {
	Username string `valid:"notnull"`
	Password string `valid:"notnull"`
}

func (a *Auth) isValid() error {
	_, err := govalidator.ValidateStruct(a)
	return err
}

func NewAuth(username, password string) (*Auth, error) {

	auth := Auth{
		Username: username,
		Password: password,
	}

	err := auth.isValid()
	if err != nil {
		return nil, err
	}

	return &auth, nil
}
