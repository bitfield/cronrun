package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bitfield/cronrun"
)

const usage = `cronrun is a tool for running scheduled jobs specified by a file in crontab format.

Usage: cronrun FILE

cronrun is designed to be called from cron every minute. If you need to run cron jobs, but don't have access to the system crontabs, or you'd prefer to manage cron jobs in a file in a Git repo, for example, then cronrun is for you.

Example file format:

# Any line starting with a # character is ignored
*/5 * * * * /usr/local/bin/backup
00 01 * * * /usr/bin/security_upgrades

Running cronrun on the example file will run /usr/local/bin/backup if the current minute is divisible by 5, and will run /usr/bin/security_upgrades if the time is 01:00.

Example cron job to trigger cronrun on a given file:

* * * * * /usr/local/bin/cronrun /var/www/mysite/crontab
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	jobs, err := cronrun.JobsFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(jobs))
	refTime := time.Now()
	for _, j := range jobs {
		go func(j cronrun.Job) {
			_, err := cronrun.RunJobIfDue(j, refTime)
			if err != nil {
				log.Print(err)
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
}
