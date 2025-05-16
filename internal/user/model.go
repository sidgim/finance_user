package user

import (
	"net/http"
)

type (
	Endpoints struct {
		Create http.HandlerFunc
		Get    http.HandlerFunc
		GetAll http.HandlerFunc
		Update http.HandlerFunc
		Delete http.HandlerFunc
	}

	CreateRequest struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Phone     string `json:"phone" validate:"required,e164"`
	}

	UpdateRequest struct {
		Email string `json:"email" validate:"required,email"`
		Phone string `json:"phone" validate:"required,e164"`
	}
)
