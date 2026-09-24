package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Locker struct {
	rdb *redis.Client

	ttl time.Duration

	instanceId string
}

func NewLocker(rdb *redis.Client, ttl time.Duration, instanceId string) *Locker {
	return &Locker{
		rdb:        rdb,
		ttl:        ttl,
		instanceId: instanceId,
	}
}

func (t *Locker) tickKey(jobId string, scheduleAt time.Time) string {
	return fmt.Sprintf("tick:%s%d", jobId, scheduleAt.Unix())
}

// find out what the scheduleAt parameter is for 
func (t *Locker) Claim(ctx context.Context, jobID string, scheduleAt time.Time) (bool, error) {
	key := t.tickKey(jobID, scheduleAt)

	ok, err := t.rdb.SetNX(ctx, key, t.instanceId, t.ttl).Result()

	if err != nil {
		return false, err
	}

	return ok, nil
}
