package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// Router to handle http request
type Router struct {
	defaultMsg string
}

func (r Router) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "你好，%s", r.defaultMsg)
}

func NewDefaultRouter(rootMsg string) http.Handler {
	return Router{
		defaultMsg: rootMsg,
	}
}

type ShutdownRouter struct {
	shutdownFn func()
}

func (r ShutdownRouter) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "shutdown...")
	r.shutdownFn()
}

func NewShutdownRouter(shutdownFn func()) http.Handler {
	return ShutdownRouter{
		shutdownFn: shutdownFn,
	}
}

// Application is an util to start/stop servers
type Application struct {
	g       *errgroup.Group
	q       context.Context
	servers []*http.Server
	cancel func()

	Name    string
	Version string
}

func NewApplication(name string, version string) Application {
	g, q := errgroup.WithContext(context.Background())
	return Application{
		g:       g,
		q:       q,
		servers: make([]*http.Server, 0, 8),

		Name:    name,
		Version: version,
	}
}

func (app *Application) Serve(server *http.Server) {
	app.g.Go(func() error {
		err := server.ListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})
	app.servers = append(app.servers, server)
}

func (app *Application) ShutdownWhenSignal(sig ...os.Signal) {
	app.g.Go(func() error {
		notifyCh := make(chan os.Signal, 1)
		signal.Notify(notifyCh, sig...)

		log.Println("running in receive signal mode")

		select {
		case sign := <- notifyCh:
			log.Printf("receive signal: %s", sign)
		case <- app.q.Done():
			log.Println("quit by other")
		}

		timeoutCtx, timeoutCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer timeoutCancel()

		var err error = nil
		for _, server := range app.servers {
			err = server.Shutdown(timeoutCtx)
		}
		if err != nil {
			return err
		}
		return nil
	})
}

func (app *Application) Shutdown() {
	app.g.Go(func() error {
		var err error
		timeoutCtx, timeoutCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer timeoutCancel()
		for _, server := range app.servers {
			err = server.Shutdown(timeoutCtx)
		}
		return err
	})
}

func (app *Application) Wait() error {
	return app.g.Wait()
}

// main

func main() {
	tattyApp := NewApplication("my test app", "0.0.1")

	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      NewDefaultRouter("world"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      NewDefaultRouter("go"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server03 := &http.Server{
		Addr:         ":8082",
		Handler:      NewShutdownRouter(tattyApp.Shutdown),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	tattyApp.Serve(server01)
	tattyApp.Serve(server02)
	tattyApp.Serve(server03)

	tattyApp.ShutdownWhenSignal(syscall.SIGINT, syscall.SIGTERM)

	err := tattyApp.Wait()
	if err != nil {
		log.Fatal("app wait err: ", err)
	}
}
