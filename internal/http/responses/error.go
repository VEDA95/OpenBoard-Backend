package responses

type ErrorMessageResponse struct {
	BaseResponse
	Error GenericMessage `json:"error"`
}

type ErrorResponse[T any] struct {
	BaseResponse
	Errors T `json:"errors"`
}

type ErrorCollectionResponse[T any] struct {
	BaseCollectionResponse
	Errors []T `json:"errors"`
}

func ErrorRespMessage(code int, message string) *ErrorMessageResponse {
	return &ErrorMessageResponse{
		BaseResponse: BaseResponse{
			Code: code,
		},
		Error: GenericMessage{Message: message},
	}
}

func ErrorResp[T any](code int, data T) *ErrorResponse[T] {
	return &ErrorResponse[T]{
		BaseResponse: BaseResponse{Code: code},
		Errors:       data,
	}
}

func ErrorCollectionResp[T any](code int, errors []T) *ErrorCollectionResponse[T] {
	return &ErrorCollectionResponse[T]{
		BaseCollectionResponse: BaseCollectionResponse{
			BaseResponse: BaseResponse{
				Code: code,
			},
			Count: len(errors),
		},
		Errors: errors,
	}
}
