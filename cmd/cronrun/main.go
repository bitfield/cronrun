package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bitfield/cronrun"
)

const usage = `cronrun is a tool for running scheduled jobs specified by a single-line file in crontab format.

Usage: cronrun FILE

Example file format:

*/5 * * * * /usr/local/bin/backup

Running cronrun on the example file will run /usr/local/bin/backup if the current minute is divisible by 5.
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
