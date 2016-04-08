package redis

import (
	"fmt"
	"os"
)

type Client interface {

	// Command execute a single command and returns the result.
	Command(c Cmd, args ...interface{}) (*Reply, error, *Job)

	// Command execute multiple commands in a pipeline and returns the result.
	Commands(args ...interface{}) ([]*Reply, error, *Job)

	// Quit quits clients and cleans up.
	// it can be called more than once. the second and later calls result in
	// InfoAlreadyQuit.
	Quit() error
}

var DefaultStatusHandler func(error) = func(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}



/////////////////////////////////////////////////////////////////////

func command(jobCh chan<- *Job, cmd Cmd, args ...interface{}) (*Reply, error, *Job) {
	job, err := newSingleCmdJob(cmd, args...)
	if err != nil {
		return nil, err, nil
	}
	job.sendAndWait(jobCh)
	var reply *Reply
	if len(job.Reply) > 0 {
		reply = job.Reply[0]
	}
	return reply, job.Err, job
}

func commands(jobCh chan<- *Job, args ...interface{}) ([]*Reply, error, *Job) {
	job, err := newJob(args...)
	if err != nil {
		return nil, err, nil
	}
	job.sendAndWait(jobCh)
	return job.Reply, job.Err, job
}

func statusHandler(statusCh <-chan error, eh func(error)) {
	if eh != nil {
		for err := range statusCh {
			eh(err)
		}
	} else {
		for _ = range statusCh {
			// intentionally left blank
		}
	}
}
