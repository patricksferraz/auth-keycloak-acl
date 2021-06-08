package model

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type EmployeeClaims struct {
	ID    string   `json:"id" mapstructure:"sub" valid:"uuid"`
	Roles []string `json:"roles,omitempty" mapstructure:"roles" valid:"-"`
}

func (e *EmployeeClaims) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func NewEmployeeClaims(employeeID string, roles []string) (*EmployeeClaims, error) {

	employeeClaims := EmployeeClaims{
		ID:    employeeID,
		Roles: roles,
	}

	err := employeeClaims.isValid()
	if err != nil {
		return nil, err
	}

	return &employeeClaims, nil
}
