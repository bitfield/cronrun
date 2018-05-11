// Package cronrun parses strings in crontab format.
package cronrun

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
)

// A CronSpec represents the data parsed from a crontab line.
// `Cronexpr` is the cron expression string (e.g. `* * * * *`).
// `Command` is the remainder of the line, which cron would run as the scheduled command.
type CronSpec struct {
	Cronexpr string
	Command  string
}

// DueNow returns true if the cron expression `crontime` represents the same time as the time `now`, to the minute, and false otherwise. For example, DueNow always returns true for the cron expression `* * * *`, since that means 'run every minute'. The expression '5 * * * *' returns true if the current minute of `now` is 5. And so on.
func DueNow(crontime string, now time.Time) (bool, error) {
	expr, err := cronexpr.Parse(crontime)
	if err != nil {
		return false, fmt.Errorf("failed to parse cron expression %q: %v", crontime, err)
	}
	thisMinute := now.Truncate(time.Minute)

	// If the job is due to run now, expr.Next(now) will return the *next* time it's due.
	// So we call expr.Next with the time one second before the start of the current minute.
	// If the result is the current minute, then the job is due to run now.
	nextRunMinute := expr.Next(thisMinute.Add(-1 * time.Second))
	return thisMinute == nextRunMinute, nil
}

// SplitCrontab parses a crontab line (like `* * * * * /usr/bin/foo`) and returns a CronSpec with the `Cronexpr` and `Command` fields set to the parsed cron expression and the command, respectively.
func SplitCrontab(crontab string) (CronSpec, error) {
	fields := strings.Fields(crontab)
	if len(fields) < 6 {
		return CronSpec{}, fmt.Errorf("less than six fields in crontab %q", crontab)
	}
	cronexpr := strings.Join(fields[:5], " ")
	command := strings.Join(fields[5:], " ")
	return CronSpec{cronexpr, command}, nil
}
