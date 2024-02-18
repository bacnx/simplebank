package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bacnx/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader     = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	abort := func(ctx *gin.Context, err error) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
	}

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeader)

		fields := strings.Fields(authHeader)

		if len(fields) < 2 {
			abort(ctx, errors.New("invalid authorization header"))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			abort(ctx, fmt.Errorf("authorization type %s is not supported", authorizationType))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			abort(ctx, err)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
