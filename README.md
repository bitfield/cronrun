declLinks: false

# cronrun
`import "github.com/bitfield/cronrun"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package cronrun is a Go library for processing crontab files, and also a
command-line tool for running scheduled jobs specified by a file in crontab
format.

### Installation
Run the following command:


	go get github.com/bitfield/cronrun

### Usage
To run `cronrun`, run the following command:


	cronrun FILE

`cronrun` is designed to be called from cron every minute. If you need to run
cron jobs, but don't have access to the system crontabs, or you'd prefer to
manage cron jobs in a file in a Git repo, for example, then `cronrun` is for
you.

### File format
`cronrun` understands any valid `crontab`-format file.

Example file format:


	# Any line starting with a # character is ignored
	*/5 * * * * /usr/local/bin/backup
	00 01 * * * /usr/bin/security_upgrades
	# Blank lines or lines containing only whitespace are ignored:
	
	* * * * * /bin/echo This will run every minute!

Running `cronrun` on the example file will run `/usr/local/bin/backup` if the
current minute is divisible by 5, and will run `/usr/bin/security_upgrades` if
the time is `01:00`.

### Running from cron
Run `cronrun` from your system crontab, as you would any normal cron job.

Example cron job to run `cronrun` on a given file:


	* * * * * /usr/local/bin/cronrun /var/www/mysite/crontab




## <a name="pkg-index">Index</a>
* [func RunJobIfDue(j Job, t time.Time) (bool, error)](#RunJobIfDue)
* [type Job](#Job)
  * [func JobsFromFile(filename string) (jobs []Job, err error)](#JobsFromFile)
  * [func NewJob(crontab string) (Job, error)](#NewJob)
  * [func (job *Job) DueAt(t time.Time) (bool, error)](#Job.DueAt)
  * [func (job *Job) Run() error](#Job.Run)


#### <a name="pkg-files">Package files</a>
[cronrun.go](cronrun.go) [doc.go](doc.go) 





## <a name="RunJobIfDue">func</a> [RunJobIfDue](cronrun.go?s=3144:3194#L99)
``` go
func RunJobIfDue(j Job, t time.Time) (bool, error)
```
RunJobIfDue runs the specified job if it is due at the specified time (as
determined by DueAt), and returns true if it was in fact run, or false
otherwise, and an error if the job failed to either parse or run
successfully.




## <a name="Job">type</a> [Job](cronrun.go?s=317:368#L17)
``` go
type Job struct {
    Due     string
    Command string
}

```
A Job represents the data parsed from a crontab line. `Due` is the cron
time specification (e.g. `* * * * *`). `Command` is the remainder of the line,
which will be run as the scheduled command.







### <a name="JobsFromFile">func</a> [JobsFromFile](cronrun.go?s=1037:1095#L38)
``` go
func JobsFromFile(filename string) (jobs []Job, err error)
```
JobsFromFile reads a multi-line crontab file, ignoring comments and blank
lines or lines containing only whitespace, and returns the corresponding
slice of Jobs, or an error.


### <a name="NewJob">func</a> [NewJob](cronrun.go?s=562:602#L25)
``` go
func NewJob(crontab string) (Job, error)
```
NewJob parses a crontab line (like `* * * * * /usr/bin/foo`) and returns a
Job with the `Due` and `Command` fields set to the parsed time specification
and the command, respectively.





### <a name="Job.DueAt">func</a> (\*Job) [DueAt](cronrun.go?s=1756:1804#L66)
``` go
func (job *Job) DueAt(t time.Time) (bool, error)
```
DueAt returns true if the job would be due to run at the specified time, and
false otherwise. For example, DueAt always returns true for jobs due at `* *
* *`, since that means 'run every minute'. A job due at '5 * * * *' is DueAt
if the current minute of `t` is 5, and so on.




### <a name="Job.Run">func</a> (\*Job) [Run](cronrun.go?s=2672:2699#L86)
``` go
func (job *Job) Run() error
```
Run runs the command line specified by `job.Command`, by passing it as an
argument to "/bin/sh -c". If the command succeeds (returns zero exit status),
a nil error is returned. If the command fails (non-zero exit status), a
non-nil error containing the combined output of the command as a string is
returned.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
