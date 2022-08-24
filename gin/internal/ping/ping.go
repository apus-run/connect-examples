package ping

import (
	"context"
	"fmt"
	"github.com/apus-run/proto-def-examples/proto/ping/v1/pingv1connect"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bufbuild/connect-go"

	pingv1 "github.com/apus-run/proto-def-examples/proto/ping/v1"
)

type Service struct{}

func (s *Service) Ping(_ context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
	log.Println("Request headers:", req.Header())
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))

	res := connect.NewResponse(&pingv1.PingResponse{Data: fmt.Sprintf("Hello, %s!", req.Msg.Data)})

	res.Header().Set("Gaia-Version", "v1.0.0")

	return res, nil
}

func MainServiceTest(t *testing.T, h http.Handler) {
	t.Parallel()
	server := httptest.NewUnstartedServer(h)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
	)

	grpcClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)

	grpcWebClient := pingv1connect.NewPingServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPCWeb(),
	)

	clients := []pingv1connect.PingServiceClient{connectClient, grpcClient, grpcWebClient}
	t.Run("gaia", func(t *testing.T) {
		for _, client := range clients {
			result, err := client.Ping(context.Background(), connect.NewRequest(&pingv1.PingRequest{
				Data: "foobar",
			}))
			assert.NoError(t, err)
			assert.Equal(t, "Hello, foobar!", result.Msg.Data)
		}
	})
}
