package queue

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

// WorkerQueue is the global queue of Workers
var WorkerQueue chan chan Job

// WorkQueue is the global queue of work to dispatch
var JobQueue = make(chan Job, 250)

// NewDispatcher creates, and returns a new Dispatcher object.
func NewDispatcher(workCh chan Job, m int, d string, w string) *Dispatcher {
	workerPool := make(chan chan *Worker, m)

	return &Dispatcher{
		MaxWorkers:   m,
		WorkersDir:   d,
		WorkspaceDir: w,
		workerPoolCh: workerPool,
		Work:         make(map[uint64]*Job),
	}
}

// StartDispatcher does stuff
func (d *Dispatcher) Start() {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan Job, d.MaxWorkers)

	log.Info("Starting dispatcher...")

	// Now, create all of our workers.
	for i := 0; i < d.MaxWorkers; i++ {
		log.Infof("Dispatching Worker %d", i+1)
		wn := strconv.Itoa(i + 1)
		worker := NewWorker(i+1, d.WorkersDir+"_"+wn, d.WorkspaceDir+"_"+wn, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-JobQueue:
				log.Debug("Received job requeust from webserver...")
				go func() {
					worker := <-WorkerQueue

					log.Debug("Dispatching job request to worker...")

					worker <- work
				}()
			}
		}
	}()
}

func (d *Dispatcher) Get(i uint64) (*Job, bool) {
	job, ok := d.Work[i]
	return job, ok
}

func (d *Dispatcher) Stop(i uint64) bool {
	job, found := d.Get(i)
	if !found {
		log.Warnf("Cannot find job request to stop with id %d", i)
		return false
	}
	job.Stop()
	return true
}
