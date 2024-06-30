package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TypeSendVerifyEmail = "type:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeSendVerifyEmail, payloadJson, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("cannot enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("ID", info.ID).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).
		Msg("enqueued send verify email")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		// }
		return fmt.Errorf("failed to get user from db: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verifyEmail in db: %w", err)
	}

	subject := "Wellcome to Simple Bank!"
	verifyUrl := fmt.Sprintf("http://simplebank.example.com/verify?id=%d&secretCode=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`
    <p>Thank you <b>%s</b> for registering with us,</br></p>
    <p>Click <a href='%s'>here</a> to verify your email.</p>
  `, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(to, subject, content, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed send verify email")
	return nil
}
