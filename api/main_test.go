package api

import (
	"testing"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, borrowerService service.BorrowerService) *Server {
	server, err := NewServer(borrowerService)
	require.NoError(t, err)

	return server
}
