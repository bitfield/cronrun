package cronrun

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
)

// A Job represents the data parsed from a crontab line. `Due` is the cron
// time specification (e.g. `* * * * *`). `Command` is the remainder of the line,
// which will be run as the scheduled command.
type Job struct {
	Due     string
	Command string
}

// NewJob parses a crontab line (like `* * * * * /usr/bin/foo`) and returns a
// Job with the `Due` and `Command` fields set to the parsed time specification
// and the command, respectively.
func NewJob(crontab string) (Job, error) {
	fields := strings.Fields(crontab)
	if len(fields) < 6 {
		return Job{}, fmt.Errorf("less than six fields in crontab %q", crontab)
	}
	due := strings.Join(fields[:5], " ")
	command := strings.Join(fields[5:], " ")
	return Job{due, command}, nil
}

// JobsFromFile reads a multi-line crontab file, ignoring comments and blank
// lines or lines containing only whitespace, and returns the corresponding
// slice of Jobs, or an error.
func JobsFromFile(filename string) (jobs []Job, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		if line == "" {
			continue
		}
		j, err := NewJob(s.Text())
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// DueAt returns true if the job would be due to run at the specified time, and
// false otherwise. For example, DueAt always returns true for jobs due at `* *
// * *`, since that means 'run every minute'. A job due at '5 * * * *' is DueAt
// if the current minute of `t` is 5, and so on.
func (job *Job) DueAt(t time.Time) (bool, error) {
	expr, err := cronexpr.Parse(job.Due)
	if err != nil {
		return false, fmt.Errorf("failed to parse cron expression %q: %v", job.Due, err)
	}
	thisMinute := t.Truncate(time.Minute)

	// If the job is due to run now, expr.Next(now) will return the *next*
	// time it's due. So we call expr.Next with the time one second before
	// the start of the current minute. If the result is the current minute,
	// then the job is due to run now.
	nextRunMinute := expr.Next(thisMinute.Add(-1 * time.Second))
	return thisMinute == nextRunMinute, nil
}

// Run runs the command line specified by `job.Command`, by passing it as an
// argument to "/bin/sh -c". If the command succeeds (returns zero exit status),
// a nil error is returned. If the command fails (non-zero exit status), a
// non-nil error containing the combined output of the command as a string is
// returned.
func (job *Job) Run() error {
	cmd := exec.Command("/bin/sh", "-c", job.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command %q failed: %v: output: %s", job.Command, err, output)
	}
	return nil
}

// RunJobIfDue runs the specified job if it is due at the specified time (as
// determined by DueAt), and returns true if it was in fact run, or false
// otherwise, and an error if the job failed to either parse or run
// successfully.
func RunJobIfDue(j Job, t time.Time) (bool, error) {
	due, err := j.DueAt(t)
	if err != nil {
		return false, err
	}
	if !due {
		return false, nil
	}
	return true, j.Run()
}
