package main

import (
	"log"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	hs := http.NewServer()
	gs := grpc.NewServer()

	app := kratos.New(
		kratos.Name("demo"),
		kratos.Version("v0.0.1"),
		kratos.Server(hs, gs))
	time.AfterFunc(time.Minute, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// another solution:
//https://github.com/da440dil/go-workgroup
