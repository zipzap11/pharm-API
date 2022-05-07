package response

type StdResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrResponse struct {
	Message string `json:"message"`
}
