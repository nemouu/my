package shell

import (
	"errors"
	"fmt"
)

const (
	FG = iota + 1
	BG
	ST
)

const MaxJobs = 16

type Job struct {
	pid     int
	jid     int
	state   int
	cmdline string
}

var jobs []Job
var nextjid int = 1

func addJob(pid int, state int, cmdline string) error {
	if len(jobs) >= MaxJobs {
		return errors.New("too many jobs")
	}
	newJob := Job{
		pid:     pid,
		jid:     nextjid,
		state:   state,
		cmdline: cmdline,
	}
	jobs = append(jobs, newJob)
	nextjid++
	return nil
}

func deleteJob(pid int) error {
	for i, job := range jobs {
		if job.pid == pid {
			jobs = append(jobs[:i], jobs[i+1:]...)
			return nil
		}
	}
	return errors.New("job not found")
}

func getJobByPid(pid int) (*Job, error) {
	for i := range jobs {
		if jobs[i].pid == pid {
			return &jobs[i], nil
		}
	}
	return nil, errors.New("job not found")
}

func getJobByJid(jid int) (*Job, error) {
	for i := range jobs {
		if jobs[i].jid == jid {
			return &jobs[i], nil
		}
	}
	return nil, errors.New("job not found")
}

func listJobs() {
	for i := range jobs {
		var state string
		switch jobs[i].state {
		case FG:
			state = "Foreground"
		case BG:
			state = "Running"
		case ST:
			state = "Stopped"
		}
		fmt.Printf("[%d] (%d) %s %s\n", jobs[i].jid, jobs[i].pid, state, jobs[i].cmdline)
	}
}
