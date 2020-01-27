package server

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

// WorkerQueue defines a worker queue.
var WorkerQueue chan chan WorkRequest

// StartDispatcher is a dispatcher
func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	sworkers := strconv.Itoa(nworkers)

	log.Info("Starting workers... (" + sworkers + " )")

	// Start the workers
	for i := 0; i < nworkers; i++ {
		log.Debug("Starting worker >> ", i+1)

		worker := NewWorker(i+1, WorkerQueue)

		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				log.Info("Received work request...")
				go func() {
					worker := <-WorkerQueue

					log.Info("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}
