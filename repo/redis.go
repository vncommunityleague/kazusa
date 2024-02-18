package repo

import (
	"github.com/redis/rueidis"
	"os"
)

func ConnectToRedis() (rueidis.Client, error) {
	url := os.Getenv("REDIS_URL")

	client, err := rueidis.NewClient(rueidis.MustParseURL(url))
	if err != nil {
		return nil, err
	}

	return client, nil
}
