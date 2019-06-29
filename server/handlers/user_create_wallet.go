package handlers

import validation "github.com/go-ozzo/ozzo-validation"

type UserWalletsRequest struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func (u UserWalletsRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(u.Name, validation.Required),
		validation.Field(u.Kind, validation.Required),
	)
}
