package worker

import (
	"strings"

	"github.com/hibiken/asynq"
)

func NewAsynqClient(option *asynq.RedisClientOpt) *asynq.Client {
	redisUrl := option.Addr
	redisUrl = strings.TrimPrefix(redisUrl, "redis://")
	return asynq.NewClient(
		asynq.RedisClientOpt{Addr: redisUrl, DB: option.DB},
	)
}
