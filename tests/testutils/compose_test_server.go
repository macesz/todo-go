package testutils

import (
	"net/http/httptest"
	"testing"

	"github.com/macesz/todo-go/cmd/composition"
	"github.com/macesz/todo-go/delivery/web"
	"github.com/macesz/todo-go/domain"
)

func ComposeServer(t *testing.T) (*TestContainer, *httptest.Server, *web.ServerServices) {
	ctx := t.Context()
	cfg := domain.Config{
		JWTSecret: "my-super-secret-test-key-12345",
	}

	// Setup database
	tc := SetupTestDB(t)

	services := composition.ComposeServices(cfg, tc.DB)

	handlers, err := web.CreateHandlers(ctx, services)
	if err != nil {
		t.Error(err)
	}

	router, err := web.CreateRouter(ctx, cfg, services, handlers)
	if err != nil {
		t.Error(err)
	}

	server := httptest.NewServer(router)

	return tc, server, services
}
