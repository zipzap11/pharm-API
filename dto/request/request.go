package request

import "github.com/zipzap11/pharm-API/model"

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required,numeric"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ModelFromCreateUserRequest(req *CreateUserRequest) *model.User {
	return &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	}
}

type RefreshSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
}
