package response

import "github.com/zipzap11/pharm-API/model"

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
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func ModelToCategoryResponse(model *model.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:       int64(model.ID),
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
