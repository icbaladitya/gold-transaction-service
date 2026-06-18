package domain

type BasicResponse[T any] struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
	Data       *T     `json:"data,omitempty"`
	Items      []T    `json:"items,omitempty"`
}

func SuccessDataResponse[T any](data T, message string) BasicResponse[T] {
	return BasicResponse[T]{
		ResultCode: 1,
		Message:    message,
		Data:       &data,
	}
}

func SuccessListResponse[T any](items []T, message string) BasicResponse[T] {
	return BasicResponse[T]{
		ResultCode: 1,
		Message:    message,
		Items:      items,
	}
}

func SuccessResponse[T any](message string) BasicResponse[T] {
	return BasicResponse[T]{
		ResultCode: 1,
		Message:    message,
	}
}

func FailResponse[T any](message string) BasicResponse[T] {
	return BasicResponse[T]{
		ResultCode: 0,
		Message:    message,
	}
}

func ErrorResponse[T any](message string) BasicResponse[T] {
	return BasicResponse[T]{
		ResultCode: -1,
		Message:    message,
	}
}
