package reviewgoroutineerrors

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

func fetchData(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", errors.New("empty url")
	}
	return "data from " + url, nil
}

func fetchAll(urls []string) ([]string, error) {
	results := make([]string, len(urls))
	errCh := make(chan error, len(urls))
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			data, err := fetchData(context.Background(), url)
			if err != nil {
				errCh <- err
				return
			}
			results[i] = data
		}(i, url)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err // devolve o primeiro erro encontrado
		}
	}

	return results, nil
}

func fetchAllWithErrgroup(urls []string) ([]string, error) {
	results := make([]string, len(urls))
	g, ctx := errgroup.WithContext(context.Background())

	for i, url := range urls {
		i, url := i, url // necessário em versões antigas do Go — hoje, opcional, mas comum ver em código real
		g.Go(func() error {
			data, err := fetchData(ctx, url)
			if err != nil {
				return err
			}
			results[i] = data
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
