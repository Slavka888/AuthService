package suite

import (
	"AuthService/internal/config"
	"context"
	"net"
	authv1 "protos/gen/go/auth-service"
	"strconv"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg         *config.Config
	AuthService authv1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadTest()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}

	return ctx, &Suite{
		T:           t,
		Cfg:         cfg,
		AuthService: authv1.NewAuthClient(cc),
	}
}
