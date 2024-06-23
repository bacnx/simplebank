package worker

import (
	"context"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCretical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(context.Context, *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisProcessor(redisOpt asynq.RedisConnOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCretical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("Faild to process task")
		}),
		Logger: NewLogger(),
	})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TypeSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
