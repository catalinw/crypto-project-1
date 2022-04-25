package public

type ApiResponse struct {
	Result  interface{} `json:"result"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
}
