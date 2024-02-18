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
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func createAndSetAuthHeader(t *testing.T, tokenMaker token.Maker, username string, request *http.Request) {
	if username == "" {
		return
	}

	setAuthHeader(t, username, time.Minute, authorizationTypeBearer, tokenMaker, request)
}
