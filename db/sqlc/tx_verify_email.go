package db

import (
	"context"
	"database/sql"
	"errors"
)

type VerifyEmailTxParam struct {
	EmailId    int64  `json:"email_id"`
	SecretCode string `json:"secret_code"`
}

type VerifyEmailTxResult struct {
	VerifyEmail VerifyEmail `json:"verify_email"`
	User        User        `json:"user"`
}

func (store *SQLStore) VerifyEmailTx(
	ctx context.Context,
	param VerifyEmailTxParam,
) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         param.EmailId,
			SecretCode: param.SecretCode,
		})
		if err != nil {
			return errors.New("cannot update verify email")
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		return err
	})

	return result, err
}
