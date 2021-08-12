package entity

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type PasswordInfo struct {
	UserID    string
	Password  string
	Temporary bool
}

func NewPasswordInfo(userID, password string, temporary bool) *PasswordInfo {
	return &PasswordInfo{
		UserID:    userID,
		Password:  password,
		Temporary: temporary,
	}
}

type User struct {
	ID         string `json:"id" valid:"uuid"`
	Username   string `json:"username,omitempty" valid:"required"`
	EmployeeID string `json:"employee_id,omitempty" attr:"employee_id" valid:"-"`
}

func (e *User) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func NewUser(username, employeeID string) (*User, error) {

	user := User{
		Username:   username,
		EmployeeID: employeeID,
	}
	user.ID = uuid.NewV4().String()

	if err := user.isValid(); err != nil {
		return nil, err
	}

	return &user, nil
}
