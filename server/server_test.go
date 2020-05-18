package server

import (
	"os"
	"testing"
	"time"
)

func TestServerShutdown(t *testing.T) {
	config := Config{
		LogLvl: "DEBUG",
		Port:   "0",
		TLS:    false,
		Cert:   "",
		Key:    "",
	}

	go func() {
		Start(config)
	}()

	time.Sleep(2 * time.Second)

	stop <- os.Interrupt
}
