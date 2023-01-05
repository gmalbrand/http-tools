package utils

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func InitLog() {
	switch {
	case os.Getenv("DEBUG") != "":
		log.SetLevel(log.DebugLevel)
	case os.Getenv("WARN") != "":
		log.SetLevel(log.WarnLevel)
	case os.Getenv("ERROR") != "":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func WaitSigInterrupt(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		cancel()
	}()
}
