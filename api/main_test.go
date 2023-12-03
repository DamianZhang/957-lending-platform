package api

import (
	"testing"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, borrowerService service.BorrowerService) *Server {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	server, err := NewServer(config, borrowerService)
	require.NoError(t, err)

	return server
}
