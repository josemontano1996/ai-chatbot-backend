package controller

type errorResponse struct {
	Error string `json:"error" validate:"required"`
}

func sendErrorPayload(err error) errorResponse {
	return errorResponse{Error: err.Error()}
}

type sucessResponse[T any] struct {
	Payload T `json:"payload" validate:"required"`
}

func sendSuccessPayload[T any](payload T) sucessResponse[T] {
	return sucessResponse[T]{
		Payload: payload,
	}
}
