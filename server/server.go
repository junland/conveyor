package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

// Config struct provides configuration fields for the server.
type Config struct {
	LogLvl       string
	Access       bool
	Port         string
	PID          string
	TLS          bool
	Cert         string
	Key          string
	WorkspaceDir string
	Workers      int
}

var stop = make(chan os.Signal)

// WorkQueue is a buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 500)

// Start sets up and starts the main server application
func Start(c Config) error {

	// Get log level environment variable.
	envLvl, err := log.ParseLevel(c.LogLvl)
	if err != nil {
		fmt.Println("Invalid log level ", envLvl)
	} else {
		// Setup logging with Logrus.
		log.SetLevel(envLvl)
	}

	if c.TLS == true {
		if c.Cert == "" || c.Key == "" {
			log.Fatal("Invalid TLS configuration, please pass a file path for both CONVEYOR_KEY and CONVEYOR_CERT")
		}
	}

	log.Info("Setting up server...")

	if _, err := os.Stat(c.WorkspaceDir); os.IsNotExist(err) {
		log.Info("Workspace directory does not exist. Creating...")
		os.Mkdir(c.WorkspaceDir, 0777)
	}

	router := RegisterRoutes()

	log.Debug("Setting up logging...")

	srv := &http.Server{Addr: ":" + c.Port, Handler: AccessLogger(router, c.Access)}

	log.Info("Starting dispatcher...")

	StartDispatcher(c.Workers)

	log.Debug("Starting server on port ", c.Port)

	go func() {
		if c.TLS == true {
			err := srv.ListenAndServeTLS(c.Cert, c.Key)
			if err != nil {
				log.Fatal("ListenAndServeTLS: ", err)
			}
		}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	log.Info("Serving on port " + c.Port + ", press CTRL + C to shutdown.")

	p := CreatePID(c.PID)

	signal.Notify(stop, os.Interrupt)

	log.Warn("After notify...")

	<-stop // wait for SIGINT

	log.Warn("Shutting down server...")

	p.RemovePID()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second) // shut down gracefully, but wait no longer than 45 seconds before halting.

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	return nil
}
