package api

import (
	"net/http"
	"os"
	"testing"
	"time"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/token"
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

func createAndSetAuthHeader(t *testing.T, tokenMaker token.Maker, username string, request *http.Request) {
	if username == "" {
		return
	}

	setAuthHeader(t, username, time.Minute, authorizationTypeBearer, tokenMaker, request)
}
