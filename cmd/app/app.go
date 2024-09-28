package main

import (
	"fmt"

	"mzhn/management/internal/app"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("cannot load env: %w", err))
	}
}

//	@title		Сервис базы знаний
//	@version	1.0
func main() {
	app, _, err := app.New()
	if err != nil {
		panic(err)
	}

	app.Run()
}
