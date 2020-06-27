package midljschema

import (
	"net/http"

	"github.com/qri-io/jsonschema"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

// DefaultBadSyntaxErrorHandler is the default invalid JSON error handler for a
// new JsonSchemaValidator instance.
//
// This function simply sets a 400 status code and passes the error up to be
// serialized by the midl.Adapter's configured error serializer.
func DefaultBadSyntaxErrorHandler(err error) midl.Response {
	return midl.MakeErrorResponse(http.StatusBadRequest, err)
}

// DefaultValidationErrorHandler is the default JSON schema validation error
// handler for a new JsonSchemaValidator instance.
//
// This function sets a 400 status code and passes the validation errors up to
// be serialized as a JSON array.
func DefaultValidationErrorHandler(errs []jsonschema.KeyError) midl.Response {
	return midl.NewResponse().
		SetCode(http.StatusBadRequest).
		SetBody(errs)
}
