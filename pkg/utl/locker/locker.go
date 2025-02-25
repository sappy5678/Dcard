package locker

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
)

func New(connectString string) (rueidislock.Locker, error) {
	client, err := rueidislock.NewLocker(
		rueidislock.LockerOption{
			ClientOption:   rueidis.ClientOption{InitAddress: []string{connectString}},
			KeyMajority:    1, // let it configable
			NoLoopTracking: true,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
