package user

import (
	"context"
	"fmt"
	"github.com/apus-run/proto-def-examples/proto/user/v1/userv1connect"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	userv1 "github.com/apus-run/proto-def-examples/proto/user/v1"
	"github.com/bufbuild/connect-go"
)

type Service struct {
	StreamDelay time.Duration
}

func (s Service) Say(
	_ context.Context,
	req *connect.Request[userv1.SayRequest],
) (*connect.Response[userv1.SayResponse], error) {
	log.Println("Request-Header: ", req.Header())
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&userv1.SayResponse{
		Sentence: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("UserSayService-Version", "v1.0.0")
	return res, nil
}

func (s Service) Introduce(
	ctx context.Context,
	req *connect.Request[userv1.IntroduceRequest],
	stream *connect.ServerStream[userv1.IntroduceResponse],
) error {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	name := req.Msg.Name
	if name == "" {
		name = "Anonymous User"
	}
	intros := []string{name + ", How are you feeling today 01 ?", name + ", How are you feeling today 02 ?"}
	var ticker *time.Ticker
	if s.StreamDelay > 0 {
		ticker = time.NewTicker(s.StreamDelay)
		defer ticker.Stop()
	}
	for _, resp := range intros {
		if ticker != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}
		}
		if err := stream.Send(&userv1.IntroduceResponse{Sentence: resp}); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) Login(
	_ context.Context,
	req *connect.Request[userv1.LoginRequest],
) (*connect.Response[userv1.LoginResponse], error) {
	log.Println("Content-Type: ", req.Header().Get("Content-Type"))
	log.Println("User-Agent: ", req.Header().Get("User-Agent"))
	res := connect.NewResponse(&userv1.LoginResponse{Token: "Bearer xxxxxxxxxx"})
	res.Header().Set("UserLoginService-Version", "v1.0.0")
	return res, nil
}

func MainServiceTest(t *testing.T, h http.Handler) {
	t.Parallel()
	server := httptest.NewUnstartedServer(h)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := userv1connect.NewUserServiceClient(
		server.Client(),
		server.URL,
	)

	grpcClient := userv1connect.NewUserServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)

	grpcWebClient := userv1connect.NewUserServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPCWeb(),
	)

	clients := []userv1connect.UserServiceClient{connectClient, grpcClient, grpcWebClient}
	t.Run("gaia", func(t *testing.T) { // nolint: paralleltest
		for _, client := range clients {
			result, err := client.Say(context.Background(), connect.NewRequest(&userv1.SayRequest{
				Name: "foobar",
			}))
			assert.NoError(t, err)
			assert.Equal(t, "Hello, foobar!", result.Msg.Sentence)
		}
	})
}
