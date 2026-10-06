package workerpool

import (
	"context"
	"sync"

	"github.com/qu0ta/go-spbu-course/internal/movie"
)

type Result struct {
	ID    int
	Movie *movie.Movie
	Err   error
}

type Pool struct {
	workers int
	client  *movie.Client
}

func New(workers int, client *movie.Client) *Pool {
	return &Pool{
		workers: workers,
		client:  client,
	}
}

func (p *Pool) Run(
	ctx context.Context,
	from int,
	to int,
) []Result {
	jobs := make(chan int)
	results := make(chan Result)

	var wg sync.WaitGroup

	// Запускаем workers.
	for i := 0; i < p.workers; i++ {
		wg.Add(1)

		go p.worker(
			ctx,
			jobs,
			results,
			&wg,
		)
	}

	// Отправляем ID в jobs.
	go func() {
		defer close(jobs)

		for id := from; id <= to; id++ {
			select {
			case jobs <- id:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Когда все workers закончат работу,
	// закрываем results.
	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []Result

	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}

func (p *Pool) worker(
	ctx context.Context,
	jobs <-chan int,
	results chan<- Result,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case id, ok := <-jobs:
			if !ok {
				return
			}

			m, err := p.client.Get(ctx, id)

			result := Result{
				ID:    id,
				Movie: m,
				Err:   err,
			}

			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}
