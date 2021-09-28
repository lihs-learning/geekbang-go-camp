package cannotexit

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

// Demo 1

func leak() {
	// No writer? What the hell!?
	ch := make(chan int)

	go func() {
		val := <-ch // a goroutine will always block in here
		fmt.Println("We received a value", val)
	}()
}

// Demo 2

func search(term string) (string, error) {
	time.Sleep(5 * time.Minute)
	return "Some Values", nil
}

func processTooMuchTime(term string) error {
	record, err := search(term)
	if err != nil {
		return err
	}

	fmt.Println("Received:", record)
	return nil
}

type result struct {
	record string
	err error
}

func processLeakGoroutine(term string) error {
	resCh := make(chan result)

	go func() {
		record, err := search(term)
		resCh <- result{record, err}
	}()

	ctx := context.TODO()
	context.WithTimeout(ctx, 3 * time.Minute)

	// if select ctx.Done, no one can read from resCh
	// the goroutine above wile leak
	select {
	case <-ctx.Done():
		return errors.New("search timeout, canceled")
	case res := <- resCh:
		if res.err != nil {
			return res.err
		}
		fmt.Println("Received:", res.record)
		return nil
	}
}
