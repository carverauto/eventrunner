package eventingest

import (
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"gofr.dev/pkg/gofr"
)

// HTTPServer is an HTTP server for handling event requests.
type HTTPServer struct {
	app            *gofr.App
	eventForwarder EventForwarder
}

// NewHTTPServer creates a new HTTP server for handling event requests.
func NewHTTPServer(app *gofr.App, forwarder EventForwarder) *HTTPServer {
	return &HTTPServer{
		app:            app,
		eventForwarder: forwarder,
	}
}

// HandleEvent handles an event request, it accepts a CustomContext and returns an interface and an error.
func (s *HTTPServer) HandleEvent(cc *customctx.CustomContext) (interface{}, error) {
	tenantID, ok := cc.GetUUIDClaim("tenant_id")
	if !ok {
		return nil, NewAuthError("Missing tenant ID")
	}

	customerID, ok := cc.GetUUIDClaim("customer_id")
	if !ok {
		return nil, NewAuthError("Missing customer ID")
	}

	var eventData []byte
	if err := cc.Bind(&eventData); err != nil {
		return nil, NewProcessingError("Invalid request body")
	}

	if err := s.eventForwarder.ForwardEvent(cc.Context, tenantID, customerID, eventData); err != nil {
		return nil, NewProcessingError("Failed to forward event")
	}

	return map[string]string{"status": "received"}, nil
}
