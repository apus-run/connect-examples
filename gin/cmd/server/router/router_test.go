package router

import (
	"testing"

	"github.com/apus-run/connect-examples/gin/internal/ping"
	"github.com/apus-run/connect-examples/gin/internal/user"
)

func TestUserService(t *testing.T) {
	user.MainServiceTest(t, New())
}

func TestPingService(t *testing.T) {
	ping.MainServiceTest(t, New())
}
