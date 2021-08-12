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
	ID            string `json:"id" valid:"uuid"`
	Username      string `json:"username,omitempty" valid:"required"`
	FirstName     string `json:"first_name,omitempty" valid:"required"`
	LastName      string `json:"last_name,omitempty" valid:"required"`
	Email         string `json:"email,omitempty" valid:"email"`
	Enabled       bool   `json:"enabled,omitempty" valid:"-"`
	EmailVerified bool   `json:"email_verified,omitempty" valid:"-"`
	EmployeeID    string `json:"employee_id" attr:"employee_id" valid:"uuid"`
}

func (e *User) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func NewUser(id, username, firstName, lastName, email string, enabled, emailVerified bool, employeeID string) (*User, error) {

	user := User{
		Username:      username,
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Enabled:       enabled,
		EmailVerified: emailVerified,
		EmployeeID:    employeeID,
	}

	if id == "" {
		user.ID = uuid.NewV4().String()
	} else {
		user.ID = id
	}

	if err := user.isValid(); err != nil {
		return nil, err
	}

	return &user, nil
}
