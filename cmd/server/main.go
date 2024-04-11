package main

import "github.com/ilya372317/pass-keeper/internal/server/app"

func main() {
	a, err := app.New("config/config.yaml")
	if err != nil {
		panic(err)
	}
	if err := a.StartGRPCServer(); err != nil {
		panic(err)
	}
}
