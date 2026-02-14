package responses

type OkResponse[T any] struct {
	BaseResponse
	Data T `json:"data"`
}

type OkCollectionResponse[T any] struct {
	BaseCollectionResponse
	Data []T `json:"data"`
}

type SuccessResponse[T any] struct {
	OkResponse[T]
	GenericMessage
}

type SuccessCollectionResponse[T any] struct {
	OkCollectionResponse[T]
	GenericMessage
}

func OKResponse[T any](code int, data T) *OkResponse[T] {
	return &OkResponse[T]{
		BaseResponse: BaseResponse{
			Code: code,
		},
		Data: data,
	}
}

func OKCollectionResponse[T any](code int, data []T) *OkCollectionResponse[T] {
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

func CreateSuccessResponse[T any](code int, message string, data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		GenericMessage: GenericMessage{Message: message},
		OkResponse: OkResponse[T]{
			BaseResponse: BaseResponse{Code: code},
			Data:         data,
		},
	}
}

func CreateSuccessCollectionResponse[T any](code int, message string, data []T) *SuccessCollectionResponse[T] {
	return &SuccessCollectionResponse[T]{
		GenericMessage: GenericMessage{Message: message},
		OkCollectionResponse: OkCollectionResponse[T]{
			BaseCollectionResponse: BaseCollectionResponse{
				BaseResponse: BaseResponse{Code: code},
				Count:        len(data),
			},
			Data: data,
		},
	}
}
