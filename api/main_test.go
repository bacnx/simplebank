package api

import (
	"os"
	"testing"
	"time"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/util"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(store db.Store) *Server {
	config := util.Config{
		AccessTokenDuration: time.Minute,
	}

	return NewServer(config, store)
}
