package api

import (
	"testing"
	"time"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/util"
)

func newTestServer(t *testing.T, borrowerService service.BorrowerService) *Server {
	config := util.Config{
		RefreshTokenDuration: time.Minute,
	}

	server := NewServer(config, borrowerService)

	return server
}
