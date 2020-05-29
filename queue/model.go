package queue

type WorkStatus int

const (
	// Idle means work is not yet started.
	Idle WorkStatus = iota
	// Running means work is running.
	Running
	// Success means work over done without any errors. (Hopefully)
	Success
	// Error means work is finished but caught some errors. (Yikes)
	Error
)

type Dispatcher struct {
	MaxWorkers   int
	WorkersDir   string
	WorkspaceDir string
	Work         map[uint64]*Job
	workerPoolCh chan chan *Worker
	workQueueCh  chan Job
}

// Worker struct
type Worker struct {
	ID           int
	WorkersDir   string
	WorkspaceDir string

	workCh  chan Job
	queueCh chan chan Job
	quitCh  chan bool
}

// Job struct is a request of work for a worker
type Job struct {
	ID      uint64
	Script  string
	CmdList []string
	Status  WorkStatus
	LogFile string
	Err     error

	quitCh chan bool
}
