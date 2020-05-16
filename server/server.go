package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"strconv"

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
	WorkersDir   string
}

var stop = make(chan os.Signal)

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

	var w int
	
	w = 1
	for w <= c.Workers {
		ws := strconv.Itoa(w)
        if _, err := os.Stat(c.WorkspaceDir + "_" + ws); os.IsNotExist(err) {
			log.Info("Workspace directory does not exist. Creating...")
			os.Mkdir(c.WorkspaceDir, 0777)
			log.Debug("Created " + c.WorkspaceDir + "_" + ws)
		}
        w = w + 1
	}
	
	w = 1
	for w <= c.Workers {
		ws := strconv.Itoa(w)
        if _, err := os.Stat(c.WorkersDir + "_" + ws); os.IsNotExist(err) {
			log.Info("Worker directory does not exist. Creating...")
			os.Mkdir(c.WorkersDir, 0777)
			log.Debug("Created " + c.WorkersDir + "_" + ws)
		}
        w = w + 1
    }

	router := c.RegisterRoutes()

	log.Debug("Setting up http logging...")

	srv := &http.Server{Addr: ":" + c.Port, Handler: AccessLogger(router, c.Access)}

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

	p := CreatePID(c.PID)

	signal.Notify(stop, os.Interrupt)

	log.Info("Serving on port " + c.Port + ", press CTRL + C to shutdown.")

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
