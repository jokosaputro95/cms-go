package api

// Response adalah struct standar untuk format respons API
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// NewSuccessResponse membuat respons sukses
func NewSuccessResponse(message string, data interface{}, meta interface{}) *Response {
	return &Response{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

// NewErrorResponse membuat respons error
func NewErrorResponse(message string) *Response {
	return &Response{
		Status:  false,
		Message: message,
	}
}