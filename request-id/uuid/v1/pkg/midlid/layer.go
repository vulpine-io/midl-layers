package midlid

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

// KeyRequestId contains the map key used to store the generated request ID in
// the request's additional context.
const KeyRequestId = "request-id"

// NewRequestIdProvider returns a new midl.Middleware instance that generates
// a unique id for a request using google/uuid and appends it to the request's
// additional context as a string under the name defined in `KeyRequestId`.
func NewRequestIdProvider() midl.Middleware {
	return new(requestIdLayer)
}

type requestIdLayer struct {}

func (r *requestIdLayer) Handle(q midl.Request) midl.Response {
	id, err := uuid.NewRandom()
	if err != nil {
		return midl.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetError(err)
	}
	q.AdditionalContext()[KeyRequestId] = id.String()
	return nil
}

