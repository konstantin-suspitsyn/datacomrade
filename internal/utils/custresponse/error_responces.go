package custresponse

import (
	"fmt"
	"net/http"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/jsonlog"
)

func LogError(r *http.Request, err error) {
	jsonlog.LogResponseErrorNoCustomProperties(r, err)
}

// ErrorResponse() is generic response for errors to the client with status code and JSON message payload
// message is an any data we want to send to user
// message_caption is a header for message to generate JSON with structure: {message_caption: message}
func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	responceJSON := make(map[string]any)

	errorJSON := make(map[string]any)

	errorJSON["message"] = message

	responceJSON["error"] = errorJSON

	err := WriteJSON(w, status, message, nil)
	if err != nil {
		LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Sends custom 404 status
func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The resource was not found"
	ErrorResponse(w, r, http.StatusNotFound, message)
}

// Sends custom 405 status
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The method %s is not allowed", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// Any unexpected error server encounters
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)

	message := "Server got a problem, serving your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

// Sends 400 status
func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := make(map[string]string)
	message["message"] = err.Error()
	ErrorResponse(w, r, http.StatusBadRequest, message)
}

// Sends 422 response. Will be used when JSON validation is failed
func FailedValidationResponse(w http.ResponseWriter, r *http.Request, err error, errorMaps map[string]string) {
	errorJson := make(map[string]any)
	errorJson["message"] = err.Error()
	errorJson["errors"] = errorMaps
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errorJson)
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	errorJson := make(map[string]any)
	errorString := "Unauthorized"
	errorJson["message"] = errorString
	ErrorResponse(w, r, http.StatusUnauthorized, errorJson)

}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusUnauthorized, "passwords do not match")
}
