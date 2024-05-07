package gapi

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/pb"
	"github.com/bacnx/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "username not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err.Error())
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed create refresh token: %s", err.Error())
	}

	createSessionParams := db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}
	_, err = server.store.CreateSession(ctx, createSessionParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %s", err.Error())
	}

	res := &pb.LoginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		User:                  convertUser(user),
	}
	return res, nil
}
