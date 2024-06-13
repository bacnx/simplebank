package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(context.Context, *PayloadSendVerifyEmail, ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisDistributor(opt asynq.RedisConnOpt) TaskDistributor {
	client := asynq.NewClient(opt)
	return &RedisTaskDistributor{
		client: client,
	}
}
