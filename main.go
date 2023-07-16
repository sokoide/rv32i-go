package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.Info("* Started")

	log.Info("* Completed")
}
