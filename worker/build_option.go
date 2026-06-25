package worker

import (
	"log"
	"net/url"

	"github.com/hibiken/asynq"
)

func BuildClientOptionFromURL(rawURL string) *asynq.RedisClientOpt {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}
	password, _ := parsedURL.User.Password()
	return &asynq.RedisClientOpt{
		Addr:     parsedURL.Host,
		Username: parsedURL.User.Username(),
		Password: password,
	}

}
