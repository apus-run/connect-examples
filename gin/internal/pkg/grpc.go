package pkg

import (
	"net/http"

	"github.com/bufbuild/connect-go"
)

var Compress1KB = connect.WithCompressMinBytes(1024)

// RouteFn gRPC route registration
type RouteFn func() (string, http.Handler)
