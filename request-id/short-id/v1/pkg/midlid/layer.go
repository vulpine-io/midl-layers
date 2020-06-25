package midlid

import (
	"net/http"

	"github.com/teris-io/shortid"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

// KeyRequestId contains the map key used to store the generated request ID in
// the request's additional context.
const KeyRequestId = "request-id"

// NewRequestIdProvider returns a new midl.Middleware instance that generates
// a unique id for a request using teris-id/shortid and appends it to the
// request's additional context under the name defined in `KeyRequestId`.
func NewRequestIdProvider() midl.Middleware {
	return new(requestIdLayer)
}

type requestIdLayer struct {}

func (r *requestIdLayer) Handle(q midl.Request) midl.Response {
	id, err := shortid.Generate()
	if err != nil {
		return midl.NewResponse().
			SetCode(http.StatusInternalServerError).
			SetError(err)
	}
	q.AdditionalContext()[KeyRequestId] = id
	return nil
}

