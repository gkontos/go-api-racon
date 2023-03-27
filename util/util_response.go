package util

import (
	"encoding/json"
	"net/http"

	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
)

// ReturnErrorJSON will create and return a json encoded exception
func ReturnErrorJSONWithCode(w http.ResponseWriter, err error, httpStatus int) {
	type errorResponse struct {
		Message string `json:"error"`
	}
	response := errorResponse{
		Message: err.Error(),
	}
	w.WriteHeader(httpStatus)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logger.Logger.Error().Err(encodeErr).Msg("Error encoding JSON for response")
	}
}
func ReturnErrorJSON(w http.ResponseWriter, err error) {

	httpStatus := http.StatusInternalServerError
	if err != nil {
		switch err.(type) {
		case *RequestParseError:
			httpStatus = http.StatusBadRequest
		case *model.ValidationError:
			httpStatus = http.StatusBadRequest // 422 // http.StatusUnprocessableEntity -- appengine does not like this httpStatus
		case *model.ResourceDoesNotExistError:
			httpStatus = http.StatusNotFound
		default:
			err = &model.GenericError{
				Err: err,
			}
			httpStatus = http.StatusInternalServerError
		}
	}
	ReturnErrorJSONWithCode(w, err, httpStatus)
}

// ReturnBodyJSON will return the body object as a JSON message
func ReturnBodyJSON(w http.ResponseWriter, body interface{}, httpStatus int) {
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Logger.Error().Err(err).Msg("Error encoding JSON for response")
	}
}

// ReturnBlankJSON will return an empty JSON response
func ReturnBlankJSON(w http.ResponseWriter, httpStatus int) {
	w.WriteHeader(httpStatus)
}
