package queue

func (j *Job) Stop() {
	j.quitCh <- true
}
