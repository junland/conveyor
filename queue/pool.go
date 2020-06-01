package queue

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/go-cmd/cmd"
	log "github.com/sirupsen/logrus"
)

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		NumWorkers: numWorkers,
		jobs:       make(map[uint64]*Job),
	}
}

func (pool *WorkerPool) RunCmd(j JobCmd) *Job {
	cmd := cmd.NewCmd("bash", "-c", fmt.Sprintf("%s %s", j.Name, strings.Join(j.Args, " ")))

	job := NewJob(j.JobId, cmd)

	pool.jobs[job.ID] = job

	go func() {
		pool.jobChan <- job
	}()

	return job
}

func (pool *WorkerPool) RunScript(j JobScriptCmd) *Job {
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	script := GenQScript(j.JobId, j.Script, j.Cmds)

	cmd := cmd.NewCmdOptions(cmdOptions, "bash", "-c", fmt.Sprintf("%s %s", script, strings.Join(j.Args, " ")))

	job := NewJob(j.JobId, cmd)

	pool.jobs[job.ID] = job

	go func() {
		pool.jobChan <- job
	}()

	return job
}

func (pool *WorkerPool) GetJob(JobID uint64) (*Job, bool) {
	job, ok := pool.jobs[JobID]
	return job, ok
}

func (pool *WorkerPool) StopJob(JobID uint64) bool {
	job, found := pool.GetJob(JobID)
	if !found {
		log.Warnf("Cannot find Job to stop with id %d", JobID)
		return false
	}
	job.Stop()
	return true
}

func (pool *WorkerPool) Start(d string, w string) (<-chan struct{}, <-chan struct{}) {
	pool.statusChan = make(chan struct{}, 1)
	pool.statusForcedChan = make(chan struct{}, 1)

	pool.jobChan = make(chan *Job, 100)
	pool.cancelChan = make(chan struct{})
	pool.forceCancelChan = make(chan struct{})

	log.Infof("Starting %d workers", pool.NumWorkers)

	for i := 0; i < pool.NumWorkers; i++ {
		pool.workersWg.Add(1)
		wn := strconv.Itoa(i + 1)
		w := NewWorker(pool, i, d+"_"+wn, w+"_"+wn)
		w.Start()
	}

	return pool.statusChan, pool.statusForcedChan
}

func (pool *WorkerPool) Stop() {
	atomic.StoreUint64(&pool.terminationFlag, 1)
	close(pool.cancelChan)

	// Wait for all workers to finish current work
	pool.workersWg.Wait()

	// Signal user we're shutting down
	close(pool.statusChan)
}

func (pool *WorkerPool) ForceStop() {
	atomic.StoreUint64(&pool.terminationFlag, 1)
	close(pool.forceCancelChan)

	pool.workersWg.Wait()

	// Signal user we're shutting down
	close(pool.statusForcedChan)
}
