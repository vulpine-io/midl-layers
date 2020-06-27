package midljschema

import (
	"encoding/json"

	"github.com/qri-io/jsonschema"
	"github.com/vulpine-io/httpx/v1/pkg/httpx/header"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

// JsonSchemaValidator defines a middleware layer that validates incoming
// JSON request bodies against a given schema.
type JsonSchemaValidator interface {
	midl.Middleware

	// BadSyntaxHandler sets the handler to be used for converting an error caused
	// by a bad json payload into an http response.
	BadSyntaxHandler(func(error) midl.Response) JsonSchemaValidator

	// ValidationHandler sets the handler to be used for converting a json schema
	// validation error into an http response.
	ValidationHandler(func([]jsonschema.KeyError) midl.Response) JsonSchemaValidator
}

// NewJsonSchemaValidator constructs a new instance of JsonSchemaValidator
// around the given schema.
//
// The constructed schema validator will have default error handlers set.
func NewJsonSchemaValidator(schema []byte) JsonSchemaValidator {
	out := new(schemaVal)

	out.schema = new(jsonschema.Schema)
	out.f400 = DefaultBadSyntaxErrorHandler
	out.f422 = DefaultValidationErrorHandler

	if err := json.Unmarshal(schema, out.schema); err != nil {
		panic(err)
	}

	return out
}

type schemaVal struct {
	schema *jsonschema.Schema
	f500   func(error) midl.Response
	f422   func([]jsonschema.KeyError) midl.Response
	f400   func(error) midl.Response
}

func (s *schemaVal) Handle(q midl.Request) midl.Response {
	if val, ok := q.Header(header.ContentType); !ok || val != "application/json" {
		return nil
	}

	errs, err := s.schema.ValidateBytes(nil, q.Body())

	if err != nil {
		return s.f400(err)
	}

	if len(errs) > 0 {
		return s.f422(errs)
	}

	return nil
}

func (s *schemaVal) BadSyntaxHandler(f func(error) midl.Response) JsonSchemaValidator {
	s.f400 = f
	return s
}

func (s *schemaVal) ValidationHandler(f func([]jsonschema.KeyError) midl.Response) JsonSchemaValidator {
	s.f422 = f
	return s
}
