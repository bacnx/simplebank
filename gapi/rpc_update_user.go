package gapi

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/pb"
	"github.com/bacnx/simplebank/util"
	"github.com/bacnx/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	tokenPayload, err := server.authorizeRequest(ctx)
	if err != nil {
		return nil, unautheticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if len(violations) != 0 {
		return nil, violationsError(violations)
	}

	if tokenPayload.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err.Error())
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "username not found")
		}
		return nil, status.Errorf(codes.Internal, "cannot update user: %s", err.Error())
	}

	return &pb.UpdateUserResponse{
		User: convertUser(user),
	}, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	return
}
