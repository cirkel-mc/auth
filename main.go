package main

import (
	"comrades-mc/auth/cmd"

	"github.com/comrades-mc/goutils/config"
	_ "github.com/joho/godotenv/autoload"
)

const serviceName = "auth"

func main() {
	cfg := config.New(serviceName)
	defer cfg.Exit()

	srv := cmd.Serve(cfg)
	srv.Run()
}
