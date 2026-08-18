package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/c-dvoid/teamwork_api_go/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(w http.ResponseWriter, r *http.Request, dest any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decode json (%v): %w", err, core_errors.ErrInvalidArgument)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf("request validation (%v): %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
