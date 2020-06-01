package queue

import (
	"github.com/go-cmd/cmd"
)

func NewJob(ID uint64, command *cmd.Cmd) *Job {
	return &Job{
		ID:            ID,
		StdoutChan:    make(chan *cmd.Status),
		InterruptChan: make(chan struct{}),
		cmd:           command,
		CmdStatus: cmd.Status{
			Cmd:      "",
			PID:      0,
			Complete: false,
			Exit:     -1,
			Error:    nil,
			Runtime:  0,
		},
	}
}

func (j *Job) Stop() {
	j.InterruptChan <- struct{}{}
}
