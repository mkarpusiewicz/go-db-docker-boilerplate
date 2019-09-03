package healthcheck

import (
	"context"
	"fmt"
	"sync"

	"github.com/mkarpusiewicz/go-db-docker-boilerplate/internal/database"
)

func Perform(ctx context.Context) ([]error, bool) {
	var errors []error

	errChan := make(chan error)
	var pwg sync.WaitGroup

	pwg.Add(1)
	go func() {
		defer pwg.Done()
		if err := testDatabase(ctx); err != nil {
			errChan <- err
		}
	}()

	var cwg sync.WaitGroup

	cwg.Add(1)
	go func() {
		defer cwg.Done()
		for err := range errChan {
			errors = append(errors, err)
		}
	}()

	pwg.Wait()
	close(errChan)
	cwg.Wait()

	return errors, len(errors) == 0
}

func testDatabase(ctx context.Context) error {
	if err := database.Connection.DB().PingContext(ctx); err != nil {
		//todo change fmt to errors wrap
		return fmt.Errorf("database client error: %v", err)
	}
	return nil
}
