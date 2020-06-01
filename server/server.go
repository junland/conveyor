package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/junland/conveyor/queue"

	log "github.com/sirupsen/logrus"
)

// Config struct provides configuration fields for the server.
type Config struct {
	LogLvl       string
	Access       bool
	Port         string
	TLS          bool
	Cert         string
	Key          string
	Workers      int
	WorkersDir   string
	WorkspaceDir string
	JobDir       string
	WorkerPool   *queue.WorkerPool
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

	if _, err := os.Stat(c.JobDir); os.IsNotExist(err) {
		log.Info("Jobs directory does not exist. Creating...")
		os.MkdirAll(c.JobDir, 0600)
		log.Debug("Created " + c.JobDir)
	}

	var w int

	// Create worker history and scripts dir.
	w = 1
	for w <= c.Workers {
		ws := strconv.Itoa(w)
		if _, err := os.Stat(c.WorkersDir + "_" + ws); os.IsNotExist(err) {
			log.Info("Worker directory does not exist. Creating...")
			os.MkdirAll(c.WorkersDir+"_"+ws, 0700)
			log.Debug("Created " + c.WorkersDir + "_" + ws)
		}
	}

	// Create worker workspaces.
	w = 1
	for w <= c.Workers {
		ws := strconv.Itoa(w)
		if _, err := os.Stat(c.WorkspaceDir + "_" + ws); os.IsNotExist(err) {
			log.Info("Workspace directory does not exist. Creating...")
			os.MkdirAll(c.WorkspaceDir+"_"+ws, 0700)
			log.Debug("Created " + c.WorkspaceDir + "_" + ws)
		}
		w = w + 1
	}

	log.Info("Setting up queue engine...")

	signal.Notify(stop, os.Interrupt)

	workerpool := queue.NewWorkerPool(c.Workers)

	poolStatusChan, poolStatusForcedChan := workerpool.Start(c.WorkersDir, c.WorkspaceDir)

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

	log.Info("Serving on port " + c.Port + ", press CTRL + C to shutdown.")

	select {
	case sig := <-stop:
		log.Printf("Caught signal %+v, gracefully stopping worker pool", sig)

		go func() {
			workerpool.Stop()
		}()

		select {
		case sig := <-stop:
			log.Warnf("Caught signal %+v, forcing pool shutdown", sig)
			workerpool.ForceStop()

		case <-poolStatusChan:
			log.Warnf("Pool has been gracefully terminated")
			return nil

		case <-poolStatusForcedChan:
			log.Warnf("Pool has been forcefully terminated")
			return nil
		}
	}

	return nil
}
