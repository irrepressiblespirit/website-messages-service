package grpcapi

import (
	"encoding/json"
	"errors"

	"github.com/irrepressiblespirit/website-messages-service/pkg/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseError(err error) error {
	var cerr entity.CoreError
	var merr entity.MessageError
	if errors.As(err, &merr) {
		return getJSONError(merr)
	}
	if errors.As(err, &cerr) {
		return getJSONError(entity.MessageError{
			StatusCode: uint64(cerr.Code),
			Message:    cerr.Message,
		})
	}
	if errors.Is(err, entity.RoomNotFoundError{}) {
		return getJSONError(entity.MessageError{
			StatusCode: uint64(codes.NotFound),
			Message:    err.Error(),
		})
	}
	return getJSONError(entity.MessageError{
		StatusCode: uint64(codes.Internal),
		Message:    err.Error(),
	})
}

func getJSONError(messageError entity.MessageError) error {
	json, err := json.Marshal(messageError)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return status.Error(codes.Code(messageError.StatusCode), string(json))
}
