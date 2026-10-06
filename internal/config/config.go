package config

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	From    int
	To      int
	Workers int
	Timeout time.Duration
}

func Parse() (Config, error) {
	var (
		from         int
		to           int
		workers      int
		timeoutSecs  int
		fromProvided bool
		toProvided   bool
	)

	flag.IntVar(&from, "from", 0, "ID первого фильма")
	flag.IntVar(&to, "to", 0, "ID последнего фильма")
	flag.IntVar(&workers, "workers", 10, "количество воркеров")
	flag.IntVar(&timeoutSecs, "timeout", 5, "таймаут одного HTTP-запроса в секундах")

	flag.Parse()

	flag.CommandLine.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "from":
			fromProvided = true
		case "to":
			toProvided = true
		}
	})

	if !fromProvided {
		return Config{}, fmt.Errorf("--from является обязательным флагом\n")
	}

	if !toProvided {
		return Config{}, fmt.Errorf("--to является обязательным флагом\n")
	}

	if from < 0 {
		return Config{}, fmt.Errorf("--from должен быть неотрицательным\n")
	}

	if to < 0 {
		return Config{}, fmt.Errorf("--to должен быть неотрицательным\n")
	}

	if from > to {
		return Config{}, fmt.Errorf("--from не может быть больше --to\n")
	}

	if workers <= 0 {
		return Config{}, fmt.Errorf("--workers должен быть больше 0\n")
	}

	if timeoutSecs <= 0 {
		return Config{}, fmt.Errorf("--timeout должен быть больше 0\n")
	}

	return Config{
		From:    from,
		To:      to,
		Workers: workers,
		Timeout: time.Duration(timeoutSecs) * time.Second,
	}, nil
}
