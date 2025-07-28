package custresponse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Pager will be used for
type Pager struct {
	PageNo       int
	ItemsPerPage int
}

type envelope_json any

// Converts data to valid JSON and sends API response with JSON
// @param status - http status to send
// @param data_caption - name for JSON wrapper. { envelope_name: data }
// @param data - data to send. Will be converted by json.Marshal
func WriteJSON(w http.ResponseWriter, status int, data any, header http.Header) error {
	json_response, err := json.Marshal(data)

	if err != nil {
		return err
	}

	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Write(json_response)

	return nil
}

// Helper function to read JSON and return custom error. To hide API implementation from frontend users
// @param struct_for_json - struct that r.Body should be converted into
func ReadJSON(w http.ResponseWriter, r *http.Request, struct_for_json any) error {

	// Limit JSON size to 1Mb
	max_bytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(struct_for_json)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("Request body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Request body contains badly-formed JSON")

		// Error if json field has incorrect input type
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Request body contains incorrect type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("Request body contains incorrect JSON type at character %d", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("Body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field_name := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("Request body contains unknown key %s", field_name)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", max_bytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// We are calling Decode second time to make sure that we get only one JSON payload, since json.Decoder is a stream
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
