package queue

import (
	"sync"

	"github.com/go-cmd/cmd"
)

var workerPool WorkerPool

type WorkerPool struct {
	NumWorkers int
	jobCounter uint64

	terminationFlag  uint64
	statusChan       chan struct{}
	statusForcedChan chan struct{}

	jobChan         chan *Job
	cancelChan      chan struct{}
	forceCancelChan chan struct{}
	workersWg       sync.WaitGroup

	jobs map[uint64]*Job
}

type Worker struct {
	ID           int
	WorkersDir   string
	WorkspaceDir string
	pool         *WorkerPool
}

type Job struct {
	ID            uint64
	CmdStatus     cmd.Status
	StdoutChan    chan *cmd.Status
	InterruptChan chan struct{}
	worker        *Worker
	cmd           *cmd.Cmd
}

type JobCmd struct {
	JobId uint64
	Name  string
	Args  []string
	Env   []string
	Dir   string
}

type JobScriptCmd struct {
	JobId     uint64
	Script    string
	ScriptDir string
	Cmds      []string
	Args      []string
}
