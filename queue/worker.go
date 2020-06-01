package queue

import (
	"strconv"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewWorker(pool *WorkerPool, ID int, d string, w string) *Worker {
	return &Worker{
		pool:         pool,
		ID:           ID,
		WorkersDir:   d,
		WorkspaceDir: w,
	}
}

func (w *Worker) Start() {
	go func() {
		defer w.pool.workersWg.Done()

		log.Infof("Ready to process jobs")

		for {
			select {
			case <-w.pool.cancelChan:
				log.Infof("Received shutdown signal, won't process any new jobs")
				return

			case job := <-w.pool.jobChan:
				// if cancelChan and jobChan have messages ready at the same time, go scheduler
				// randomly selected one of the select cases. So it can happen that the job is still
				// scheduled (and if very unlucky, it can happen more than once in a row too)
				if atomic.LoadUint64(&w.pool.terminationFlag) == 1 {
					return
				}

				log.Infof("Processing job %d", job.ID)

				job.worker = w
				w.process(job)
			}
		}
	}()
}

func (w *Worker) process(j *Job) {
	j.cmd.Dir = j.worker.WorkspaceDir

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		lname := w.WorkersDir + "/job_" + strconv.FormatUint(j.ID, 10)

		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for j.cmd.Stdout != nil || j.cmd.Stderr != nil {
			select {
			case line, open := <-j.cmd.Stdout:
				if !open {
					j.cmd.Stdout = nil
					continue
				}
				AppendToFile(lname, line+"\n")
			case line, open := <-j.cmd.Stderr:
				if !open {
					j.cmd.Stderr = nil
					continue
				}
				AppendToFile(lname, line+"\n")
			}
		}
	}()

	statusChan := j.cmd.Start() // non-blocking

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			status := j.cmd.Status()
			j.CmdStatus = status

			log.Debugf("Job %d status: %v", j.ID, status)

			go func() {
				j.StdoutChan <- &status
			}()
		}
	}()

	select {
	case <-w.pool.forceCancelChan:
		log.Warnf("Forcefully stopping job %d ...", j.ID)
		j.cmd.Stop()
		ticker.Stop()
		close(j.StdoutChan)
	case <-j.InterruptChan:
		log.Infof("Requested to stop job %d", j.ID)
		ticker.Stop()
		j.cmd.Stop()
		status := j.cmd.Status()
		go func() {
			defer close(j.StdoutChan)
			j.StdoutChan <- &status
		}()
	case finalStatus := <-statusChan:
		ticker.Stop()
		j.CmdStatus = finalStatus

		go func() {
			defer close(j.StdoutChan)
			j.StdoutChan <- &finalStatus
		}()

		if !finalStatus.Complete {
			log.Warnf("Forced termination of job %d", j.ID)
			return
		}

		log.Infof("Job %d completed. Final status: %v", j.ID, finalStatus)
	}
}
