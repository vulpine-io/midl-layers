package midljschema_test

import (
	"github.com/qri-io/jsonschema"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/vulpine-io/midl/v1/pkg/midlmock"

	"github.com/vulpine-io/midl-layers/json/schema/qri-io/v1/pkg/midljschema"
)

func TestSchemaVal_Handle(t *testing.T) {
	Convey("Handle HTTP Requests", t, func() {
		Convey("Valid Request", func() {
			schema := `{"type":"array","items":{"type":"string"},"minItems":3}`
			test := midljschema.NewJsonSchemaValidator([]byte(schema))
			mock := Request{
				HeaderFunc: func(string) (string, bool) { return "application/json", true },
				BodyFunc: func() []byte {
					return []byte(`["hi", "bye", "whot"]`)
				},
			}

			So(test.Handle(&mock), ShouldBeNil)
		})

		Convey("Failing JSON", func() {
			schema := `{"type":"object","additionalProperties":{"type":"integer"}}`
			test := midljschema.NewJsonSchemaValidator([]byte(schema))
			mock := Request{
				HeaderFunc: func(string) (string, bool) { return "application/json", true },
				BodyFunc: func() []byte {
					return []byte(`{"foo": "bar"}`)
				},
			}

			res := test.Handle(&mock)
			So(res, ShouldNotBeNil)
			So(res.Code(), ShouldEqual, http.StatusBadRequest)
			So(res.Error(), ShouldBeNil)

			body := res.Body()
			So(body, ShouldNotBeNil)
			retyped, ok := body.([]jsonschema.KeyError)
			So(ok, ShouldBeTrue)
			So(len(retyped), ShouldEqual, 1)
		})

		Convey("Invalid JSON", func() {
			schema := `{"type":"object","additionalProperties":{"type":"integer"}}`
			test := midljschema.NewJsonSchemaValidator([]byte(schema))
			mock := Request{
				HeaderFunc: func(string) (string, bool) { return "application/json", true },
				BodyFunc: func() []byte {
					return []byte(`{"foo": "bar"`)
				},
			}

			res := test.Handle(&mock)
			So(res, ShouldNotBeNil)
			So(res.Code(), ShouldEqual, http.StatusBadRequest)
			So(res.Error(), ShouldNotBeNil)
		})

		Convey("Non-JSON Request", func() {
			schema := `{"type":"object","additionalProperties":{"type":"integer"}}`
			test := midljschema.NewJsonSchemaValidator([]byte(schema))
			mock := Request{
				HeaderFunc: func(string) (string, bool) { return "asdf", true },
				BodyFunc:   func() []byte { panic(nil) },
			}

			res := test.Handle(&mock)
			So(res, ShouldBeNil)
		})

	})
}
