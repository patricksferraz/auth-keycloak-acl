package entity

import (
	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Filter struct {
	Username *string `json:"username" valid:"optional"`
	Enabled  *bool   `json:"enabled" valid:"optional"`
	PageSize *int    `json:"page_size" valid:"optional"`
	Page     *int    `json:"page" valid:"optional"`
}

func (e *Filter) isValid() error {
	_, err := govalidator.ValidateStruct(e)
	return err
}

func NewFilter(username *string, enabled *bool, pageSize, page *int) (*Filter, error) {

	if pageSize == nil {
		ps := 10
		pageSize = &ps
	}

	if page == nil {
		p := 0
		page = &p
	}

	entity := &Filter{
		Username: username,
		PageSize: pageSize,
		Enabled:  enabled,
		Page:     page,
	}

	err := entity.isValid()
	if err != nil {
		return nil, err
	}

	return entity, nil
}
