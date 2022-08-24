package router

import (
	"log"
	"net/http"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/apus-run/proto-def-examples/proto/ping/v1/pingv1connect"
	"github.com/apus-run/proto-def-examples/proto/user/v1/userv1connect"

	"github.com/apus-run/connect-examples/gin/internal/grpc"
	"github.com/apus-run/connect-examples/gin/internal/pkg"
)

func grpcHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("protocol version:", c.Request.Proto)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var allServices = []string{
	userv1connect.UserServiceName,
	pingv1connect.PingServiceName,
	grpc_health_v1.Health_ServiceDesc.ServiceName,
}

func V1Route() (string, http.Handler) {
	// grpcV1
	return grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(allServices...),
		pkg.Compress1KB,
	)
}

func V1AlphaRoute() (string, http.Handler) {
	// grpcV1Alpha
	return grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(allServices...),
		pkg.Compress1KB,
	)
}

func gRPCRouter(r *gin.Engine, fn pkg.RouteFn) {
	p, h := fn()
	r.POST(p+":name", grpcHandler(h))
}

func New() *gin.Engine {
	r := gin.Default()

	gRPCRouter(r, V1Route)
	gRPCRouter(r, V1AlphaRoute)
	gRPCRouter(r, grpc.PingRoute)
	gRPCRouter(r, grpc.HealthRoute)
	gRPCRouter(r, grpc.UserRoute)

	return r
}
