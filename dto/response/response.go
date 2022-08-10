package response

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/zipzap11/pharm-API/model"
)

type StdResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrResponse struct {
	Message string `json:"message"`
}

type SessionResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TotalPriceResponse struct {
	Price         int64 `json:"price"`
	ShippingPrice int64 `json:"shipping_price"`
	TotalPrice    int64 `json:"total_price"`
}

type CategoryResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ImageURL string    `json:"image_url"`
}

func ModelToCategoryResponse(model *model.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:       model.ID,
		Name:     model.Name,
		ImageURL: model.ImageURL,
	}
}

func ModelToCategoryResponseArray(model []*model.Category) []*CategoryResponse {
	var categoryResponseArray []*CategoryResponse
	for _, category := range model {
		categoryResponseArray = append(categoryResponseArray, ModelToCategoryResponse(category))
	}
	return categoryResponseArray
}

type UserResponse struct {
	ID        uint         `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
}

func userResponseFromModel(m *model.User) *UserResponse {
	return &UserResponse{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Phone:     m.Phone,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

func UserResponsesFromModel(users []*model.User) []*UserResponse {
	var res []*UserResponse
	for _, v := range users {
		res = append(res, userResponseFromModel(v))
	}
	return res
}
