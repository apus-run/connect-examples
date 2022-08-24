package grpc

import (
	"net/http"
	"time"

	"github.com/apus-run/proto-def-examples/proto/user/v1/userv1connect"

	"github.com/apus-run/connect-examples/gin/internal/pkg"
	"github.com/apus-run/connect-examples/gin/internal/user"
)

func UserRoute() (string, http.Handler) {
	userService := &user.Service{
		StreamDelay: 2 * time.Second,
	}

	return userv1connect.NewUserServiceHandler(
		userService,
		pkg.Compress1KB,
	)
}
