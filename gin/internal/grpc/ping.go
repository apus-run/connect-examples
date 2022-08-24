package grpc

import (
	"net/http"

	"github.com/apus-run/connect-examples/gin/internal/ping"
	"github.com/apus-run/connect-examples/gin/internal/pkg"
	"github.com/apus-run/proto-def-examples/proto/ping/v1/pingv1connect"
)

func PingRoute() (string, http.Handler) {
	pingService := &ping.Service{}
	return pingv1connect.NewPingServiceHandler(pingService, pkg.Compress1KB)
}
