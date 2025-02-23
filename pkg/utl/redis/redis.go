package redis

import (
	"context"

	"github.com/redis/rueidis"
)

func New(connectString string) (rueidis.Client, error) {
	client, err := rueidis.NewClient(
		rueidis.ClientOption{
			InitAddress: []string{connectString},
		},
	)
	if err != nil {
		return nil, err
	}

	if client.Do(context.Background(), client.B().Ping().Build()).Error() != nil {
		return nil, err
	}

	return client, nil
}
