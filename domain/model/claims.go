package model

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Claims struct {
	EmployeeID string   `json:"employee_id" mapstructure:"sub" valid:"uuid"`
	Roles      []string `json:"roles,omitempty" mapstructure:"roles" valid:"-"`
}

func (e *Claims) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func NewClaims(employeeID string, roles []string) (*Claims, error) {

	claims := Claims{
		EmployeeID: employeeID,
		Roles:      roles,
	}

	err := claims.isValid()
	if err != nil {
		return nil, err
	}

	return &claims, nil
}
