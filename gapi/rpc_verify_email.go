package gapi

import (
	"context"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/pb"
	"github.com/bacnx/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) VerifyEmail(
	ctx context.Context,
	req *pb.VerifyEmailRequest,
) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, violationsError(violations)
	}

	arg := db.VerifyEmailTxParam{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	}

	_, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyEmailResponse{
		IsVerified: true,
	}, nil
}

func validateVerifyEmailRequest(
	req *pb.VerifyEmailRequest,
) (evolations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		evolations = append(evolations, fieldViolation("email_id", err))
	}
	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		evolations = append(evolations, fieldViolation("secret_code", err))
	}
	return evolations
}
