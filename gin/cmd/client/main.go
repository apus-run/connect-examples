package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	"golang.org/x/net/http2"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/apus-run/proto-def-examples/proto/ping/v1/pingv1connect"
	userv1 "github.com/apus-run/proto-def-examples/proto/user/v1"
	"github.com/apus-run/proto-def-examples/proto/user/v1/userv1connect"
)

func healthCheck(client *http.Client, services ...string) {
	healthClient := connect.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
		client,
		"http://localhost:8800/grpc.health.v1.Health/Check",
	)

	grpcHealthClient := connect.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
		client,
		"http://localhost:8800/grpc.health.v1.Health/Check",
		connect.WithGRPC(),
	)

	grpcHealthWebClient := connect.NewClient[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse](
		client,
		"http://localhost:8800/grpc.health.v1.Health/Check",
		connect.WithGRPCWeb(),
	)

	reqClients := []*connect.Client[grpc_health_v1.HealthCheckRequest, grpc_health_v1.HealthCheckResponse]{
		healthClient,
		grpcHealthClient,
		grpcHealthWebClient,
	}

	for _, n := range services {
		req := &grpc_health_v1.HealthCheckRequest{}
		if n != "" {
			req.Service = n
		}

		for _, c := range reqClients {
			res, err := c.CallUnary(
				context.Background(),
				connect.NewRequest(req),
			)
			if err != nil {
				log.Fatal(err)
			}
			if grpchealth.Status(res.Msg.Status) != grpchealth.StatusServing {
				log.Fatalf("got status %v, expected %v", res.Msg.Status, grpchealth.StatusServing)
			}
		}

	}
}

func main() {
	c := &http.Client{
		Transport: &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
			AllowHTTP: true,
		},
		Timeout: 5 * time.Second,
	}

	connectUserClient := userv1connect.NewUserServiceClient(
		c,
		"http://localhost:8800/",
	)

	grpcUserClient := userv1connect.NewUserServiceClient(
		c,
		"http://localhost:8800/",
		connect.WithGRPC(),
	)

	grpcWebUserClient := userv1connect.NewUserServiceClient(
		c,
		"http://localhost:8800/",
		connect.WithGRPCWeb(),
	)

	userClients := []userv1connect.UserServiceClient{connectUserClient, grpcUserClient, grpcWebUserClient}

	for _, client := range userClients {
		req := connect.NewRequest(&userv1.SayRequest{
			Name: "foobar",
		})
		req.Header().Set("User-Header", "hello from connect")
		res, err := client.Say(context.Background(), req)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Message:", res.Msg.Sentence)
		log.Println("UserSayService-Version:", res.Header().Get("UserSayService-Version"))
	}

	// health check
	healthCheck(c,
		userv1connect.UserServiceName,
		pingv1connect.PingServiceName,
	)
}
