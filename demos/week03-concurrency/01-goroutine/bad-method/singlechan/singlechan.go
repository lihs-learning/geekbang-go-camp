package singlechan

import (
	"fmt"
	"time"
)

type Some struct {}

// ListDirectorySlice may cost too much time
func (s Some) ListDirectorySlice(dir string) ([] string, error) {
	fmt.Println("listing", dir)
	time.Sleep(5 * time.Minute)
	return []string{"some"}, nil
}

// ListDirectoryChan may be errored but caller can't get
func (s Some) ListDirectoryChan(dir string) chan string {
	fmt.Println("listing", dir)
	ch := make(chan string)
	go func() {
		time.Sleep(5 * time.Minute)
		ch <- "some"
	}()
	return ch
}