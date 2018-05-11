// Package cronrun parses strings in crontab format.
package cronrun

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
)

// A Job represents the data parsed from a crontab line.
// `Due` is the cron expression string (e.g. `* * * * *`).
// `Command` is the remainder of the line, which cron would run as the scheduled command.
type Job struct {
	Due     string
	Command string
}

// DueNow returns true if `job` is due to run at the time specified by `now` , and false otherwise. For example, DueNow always returns true for jobs due at `* * * *`, since that means 'run every minute'. A job due at '5 * * * *' is DueNow if the current minute of `now` is 5. And so on.
func DueNow(job Job, now time.Time) (bool, error) {
	expr, err := cronexpr.Parse(job.Due)
	if err != nil {
		return false, fmt.Errorf("failed to parse cron expression %q: %v", job.Due, err)
	}
	thisMinute := now.Truncate(time.Minute)

	// If the job is due to run now, expr.Next(now) will return the *next* time it's due.
	// So we call expr.Next with the time one second before the start of the current minute.
	// If the result is the current minute, then the job is due to run now.
	nextRunMinute := expr.Next(thisMinute.Add(-1 * time.Second))
	return thisMinute == nextRunMinute, nil
}

// NewJob parses a crontab line (like `* * * * * /usr/bin/foo`) and returns a Job with the `Due` and `Command` fields set to the parsed cron expression and the command, respectively.
func NewJob(crontab string) (Job, error) {
	fields := strings.Fields(crontab)
	if len(fields) < 6 {
		return Job{}, fmt.Errorf("less than six fields in crontab %q", crontab)
	}
	due := strings.Join(fields[:5], " ")
	command := strings.Join(fields[5:], " ")
	return Job{due, command}, nil
}
