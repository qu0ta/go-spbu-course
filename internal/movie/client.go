package movie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://homeworksite.site"

type Client struct {
	httpClient *http.Client
}

type HttpError struct {
	StatusCode int
}

func (h HttpError) Error() string {
	return fmt.Sprintf("HTTP %d", h.StatusCode)
}

func FormatError(err error) string {
	if httpErr, ok := errors.AsType[*HttpError](err); ok {
		switch httpErr.StatusCode {
		case http.StatusNotFound:
			return "фильм не найден: (404)"
		case http.StatusTooManyRequests:
			return "слишком много запросов: (429)"
		case http.StatusInternalServerError:
			return "ошибка сервера: (500)"
		default:
			return fmt.Sprintf("HTTP ошибка: (%d)", httpErr.StatusCode)
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "таймаут запроса"
	}
	if errors.Is(err, context.Canceled) {
		return "запрос отменён"
	}
	return err.Error()
}

func NewClient(timeout time.Duration) *Client {
	return &Client{httpClient: &http.Client{Timeout: timeout}}
}

func (c *Client) Get(ctx context.Context, id int) (*Movie, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating an http-request: %w", err)

	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing an http-request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &HttpError{StatusCode: resp.StatusCode}
	}

	var movie Movie
	if err = json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("error while decoding JSON: %w", err)
	}
	movie.ID = id
	return &movie, nil
}
