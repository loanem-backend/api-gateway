package respx

func ResponseFail(message string, err error) *Response {
	return &Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
}

func ResponseSucceed(message string, data any) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}
