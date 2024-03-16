package model

import (
	"github.com/google/uuid"
)

// LoginAPI model
type LoginAPI struct {
	Username *string `json:"username,omitempty" example:"armadamuhammads@gmail.com" validate:"required"`
	Password *string `json:"password,omitempty" example:"password" validate:"required"`
	Remember *bool   `json:"remember,omitempty" example:"true"`
}

// RegistrationAPI model
type RegistrationAPI struct {
	FirstName       *string `json:"first_name,omitempty" example:"Armada" validate:"required"`
	LastName        *string `json:"last_name,omitempty" example:"Muhammad"`
	Email           *string `json:"email,omitempty" example:"armadamuhammads@gmail.com" validate:"required,email"`
	Username        *string `json:"username,omitempty" example:"armadamuhammads"`
	PhoneNumber     *string `json:"phone_number,omitempty" example:"089678009400"`
	Password        *string `json:"password,omitempty" example:"@Password123" validate:"required,min=8"`
	ConfirmPassword *string `json:"confirm_password,omitempty" example:"@Password123" validate:"required,eqfield=Password"`
}

type ChangePasswordAPI struct {
	Email           *string `json:"email,omitempty" example:"armadamuhammads@gmail.com" validate:"required,email"`
	Password        *string `json:"password,omitempty" example:"@Password123" validate:"required,min=8"`
	ConfirmPassword *string `json:"confirm_password,omitempty" example:"@Password123" validate:"required,eqfield=Password"`
}

// LoginResponse model
type LoginResponse struct {
	Token *ResponseToken `json:"token,omitempty"`
	User  *User          `json:"user,omitempty"`
}

// ResponseToken model
type ResponseToken struct {
	AccessToken *string `json:"access_token,omitempty"`
	ExpiresIn   *int64  `json:"expires_in,omitempty"`
	TokenType   *string `json:"token_type,omitempty"`
	IsAdmin     *bool   `json:"is_admin,omitempty"`
}

type Auth struct {
	UserID  *uuid.UUID `json:"user_id,omitempty"`
	IsAdmin *bool      `json:"is_admin,omitempty"`
	Exp     *int64     `json:"exp,omitempty"`
}
