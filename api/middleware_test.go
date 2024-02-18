package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bacnx/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setAuthHeader(
	t *testing.T,
	username string,
	duration time.Duration,
	tokenType string,
	tokenMaker token.Maker,
	request *http.Request,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", tokenType, token)
	request.Header.Set(authorizationPayloadKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setHeader     func(maker token.Maker, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setHeader: func(tokenMaker token.Maker, request *http.Request) {
				setAuthHeader(t, "user", time.Minute, authorizationTypeBearer, tokenMaker, request)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "ExpriedToken",
			setHeader: func(tokenMaker token.Maker, request *http.Request) {
				setAuthHeader(t, "user", -time.Minute, authorizationTypeBearer, tokenMaker, request)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidHeader",
			setHeader: func(tokenMaker token.Maker, request *http.Request) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "TokenTypeNotSupported",
			setHeader: func(tokenMaker token.Maker, request *http.Request) {
				setAuthHeader(t, "user", time.Minute, "NotSuported", tokenMaker, request)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			server.router.Use(authMiddleware(server.tokenMaker))

			authPath := "/auth"
			server.router.GET(authPath, func(ctx *gin.Context) {
				payload, exists := ctx.Get(authorizationPayloadKey)

				if !exists {
					ctx.JSON(http.StatusUnauthorized, nil)
					return
				}

				ctx.JSON(http.StatusOK, payload)
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setHeader(server.tokenMaker, request)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)
		})
	}
}
