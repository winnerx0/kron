package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/winnerx0/kron/internal/job"
)

type Poller struct {
	jobRepo   job.Repository
	publisher *Publisher
	locker    *Locker
	interval  time.Duration
}

func NewPoller(jobRepo job.Repository, publisher *Publisher, locker *Locker, interval time.Duration) *Poller {
	return &Poller{
		jobRepo:   jobRepo,
		publisher: publisher,
		locker:    locker,
		interval:  interval,
	}
}

func (p *Poller) Run(ctx context.Context) {

	ticker := time.NewTicker(time.Second * 5)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.tick(ctx)
		}
	}
}

func (p *Poller) tick(ctx context.Context) {

	jobs, err := p.jobRepo.FindAllNextRun(ctx)

	if err != nil {
		log.Printf("job poll failed %v", err)
		return
	}

	for _, job := range jobs {
		log.Println("job", job)

		claimed, err := p.locker.Claim(ctx, job.ID, job.NextRunAt)

		if err != nil {
			log.Println("claim failed", err)
			continue
		}

		if !claimed {
			log.Println("job already claimed")
			continue
		}

		if err := p.publisher.Publish(ctx, job); err != nil {
			log.Println("publish failed", err)
		}
	}
}
