package grpc

import (
	"net/http"

	grpchealth "github.com/bufbuild/connect-grpchealth-go"

	"github.com/apus-run/proto-def-examples/proto/ping/v1/pingv1connect"
	"github.com/apus-run/proto-def-examples/proto/user/v1/userv1connect"

	"github.com/apus-run/connect-examples/gin/internal/pkg"
)

func HealthRoute() (string, http.Handler) {
	// grpcHealthCheck
	return grpchealth.NewHandler(
		grpchealth.NewStaticChecker(
			pingv1connect.PingServiceName,
			userv1connect.UserServiceName,
		),
		pkg.Compress1KB,
	)
}
