package main

import (
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/app"
	"github.com/ilya372317/pass-keeper/pkg/logger"
)

func main() {
	logger.InitMust()
	a, err := app.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = a.ExecuteCommandCLI(); err != nil {
		fmt.Println(err)
	}
	if err = a.Stop(); err != nil {
		fmt.Println(err)
	}
}
