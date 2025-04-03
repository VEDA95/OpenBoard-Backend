package responses

type OkResponse[T interface{}] struct {
	BaseResponse
	Data T `json:"data"`
}

type OkCollectionResponse[T interface{}] struct {
	BaseCollectionResponse
	Data []T `json:"data"`
}

type SuccessResponse[T interface{}] struct {
	OkResponse[T]
	GenericMessage
}

func OKResponse[T interface{}](code int, data T) *OkResponse[T] {
	return &OkResponse[T]{
		BaseResponse: BaseResponse{
			Code: code,
		},
		Data: data,
	}
}

func OKCollectionResponse[T interface{}](code int, data []T) *OkCollectionResponse[T] {
	return &OkCollectionResponse[T]{
		BaseCollectionResponse: BaseCollectionResponse{
			BaseResponse: BaseResponse{
				Code: code,
			},
			Count: len(data),
		},
		Data: data,
	}
}

func CreateSuccessResponse[T interface{}](code int, message string, data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		GenericMessage: GenericMessage{Message: message},
		OkResponse: OkResponse[T]{
			BaseResponse: BaseResponse{Code: code},
			Data:         data,
		},
	}
}
