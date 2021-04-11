package main

import (
	"flag"

	"github.com/phsiao/fargate-migrate/internal/config"
	log "github.com/sirupsen/logrus"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "fargate-migrate config file to use")
}

func main() {
	flag.Parse()

	_, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
}
