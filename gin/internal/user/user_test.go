package user

import (
	"net/http"
	"testing"

	"github.com/apus-run/proto-def-examples/proto/user/v1/userv1connect"
)

func TestService(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(userv1connect.NewUserServiceHandler(
		&Service{},
	))
	MainServiceTest(t, mux)
}
