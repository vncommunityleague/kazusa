package repo

import (
	"os"

	"github.com/redis/rueidis"
)

func ConnectToRedis() (rueidis.Client, error) {
	url := os.Getenv("REDIS_URL")

	client, err := rueidis.NewClient(rueidis.MustParseURL(url))
	if err != nil {
		return nil, err
	}

	return client, nil
}
