package entity

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type UserEvent struct {
	ID   string `json:"id" valid:"uuid"`
	User *User  `json:"user,omitempty" valid:"-"`
}

func NewUserEvent(user *User) (*UserEvent, error) {

	e := UserEvent{
		ID:   uuid.NewV4().String(),
		User: user,
	}

	if err := e.isValid(); err != nil {
		return nil, err
	}

	return &e, nil
}

func (e *UserEvent) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func (e *UserEvent) ToJson() ([]byte, error) {
	err := e.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(e)
	if err != nil {
		return nil, nil
	}

	return result, nil
}
