package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/qu0ta/go-spbu-course/internal/config"
	"github.com/qu0ta/go-spbu-course/internal/movie"
	"github.com/qu0ta/go-spbu-course/internal/workerpool"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while parsing config: %v", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()
	client := movie.NewClient(cfg.Timeout)
	pool := workerpool.New(cfg.Workers, client)
	results := pool.Run(
		ctx, cfg.From, cfg.To)
	var movies []movie.Movie
	for _, result := range results {
		if result.Err != nil {
			fmt.Fprintf(os.Stderr, "Error movie #%d - %v\n", result.ID, movie.FormatError(result.Err))
			continue

		}
		movies = append(movies, *result.Movie)

	}
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ID < movies[j].ID
	})

	for _, m := range movies {
		fmt.Printf(
			"%d - %s - %d - %s\n",
			m.ID, m.Title, m.Year, m.Director,
		)
	}

}
