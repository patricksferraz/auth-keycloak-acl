package entity

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Claims struct {
	UserID     string   `json:"user_id" mapstructure:"sub" valid:"uuid"`
	Username   string   `json:"username" mapstructure:"preferred_username" valid:"required"`
	EmployeeID string   `json:"employee_id" mapstructure:"employee_id" valid:"uuid"`
	Roles      []string `json:"roles,omitempty" mapstructure:"roles" valid:"-"`
}

func (e *Claims) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func NewClaims(userID, username, employeeID string, roles []string) (*Claims, error) {

	claims := Claims{
		UserID:     userID,
		Username:   username,
		EmployeeID: employeeID,
		Roles:      roles,
	}

	err := claims.isValid()
	if err != nil {
		return nil, err
	}

	return &claims, nil
}
