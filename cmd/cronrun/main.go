package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	spec, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	job, err := cronrun.NewJob(string(spec))
	if err != nil {
		log.Fatal(err)
	}
	due, err := job.DueAt(time.Now())
	if err != nil {
		log.Fatal(err)
	}
	if !due {
		return
	}
	output, err := job.Run()
	if err != nil {
		log.Println(err, output)
	}
}