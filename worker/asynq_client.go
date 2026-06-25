package worker

import (
	"github.com/hibiken/asynq"
)

func NewAsynqClient(option *asynq.RedisClientOpt) *asynq.Client {
	return asynq.NewClient(option)
}
