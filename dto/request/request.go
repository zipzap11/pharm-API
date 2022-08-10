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

type UpdateItemQuantityRequest struct {
	ItemID   int64  `json:"item_id"`
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}

type CreateAddressRequest struct {
	Name          string `json:"name"`
	StateID       string `json:"state_id"`
	ProvinceID    string `json:"province_id"`
	AddressDetail string `json:"address_detail"`
	PostalCode    string `json:"postal_code"`
}

func ModelfromCreateAddressRequest(req *CreateAddressRequest) *model.Address {
	return &model.Address{
		Name:          req.Name,
		ProvinceID:    req.ProvinceID,
		StateID:       req.StateID,
		PostalCode:    req.PostalCode,
		AddressDetail: req.AddressDetail,
	}
}

type CreateTransactionRequest struct {
	AddressID        int64  `json:"address_id"`
	ShippingServices string `json:"shipping_services"`
}

type UpdateProductStockRequest struct {
	ProductID int64 `json:"product_id" validate:"required,min=1"`
	Stock     int   `json:"stock" validate:"required,min=0"`
}

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func (ccr *CreateCategoryRequest) ToModel() *model.Category {
	return &model.Category{
		Name:     ccr.Name,
		ImageURL: ccr.ImageUrl,
	}
}
