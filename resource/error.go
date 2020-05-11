package resource

import (
	"net/http"

	"github.com/ParthoShuvo/covid-19-api/errors"
)

func SendError(rw http.ResponseWriter, err error) {
	errorCode, errorText := getStatus(err)
	http.Error(rw, errorText, errorCode)
}

func getStatus(err error) (statusCode int, statusText string) {
	switch errors.GetType(err) {
	case errors.BadRequest:
		statusCode, statusText = http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
	case errors.Unauthorized:
		statusCode, statusText = http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)
	case errors.NotFound:
		statusCode, statusText = http.StatusNotFound, http.StatusText(http.StatusNotFound)
	case errors.Forbidden:
		statusCode, statusText = http.StatusForbidden, http.StatusText(http.StatusForbidden)
	default:
		statusCode, statusText = http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}
	return
}
