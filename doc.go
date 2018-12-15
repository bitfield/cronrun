// Package cronrun is a Go library for processing crontab files, and also a
// command-line tool for running scheduled jobs specified by a file in crontab
// format.
//
// Installation
//
// Run the following command:
//
// 	go get github.com/bitfield/cronrun
//
// Usage
//
// To run `cronrun`, run the following command:
//
// 	cronrun FILE
//
// `cronrun` is designed to be called from cron every minute. If you need to run
// cron jobs, but don't have access to the system crontabs, or you'd prefer to
// manage cron jobs in a file in a Git repo, for example, then `cronrun` is for
// you.
//
// File format
//
// `cronrun` understands any valid `crontab`-format file.
//
// Example file format:
//
//     # Any line starting with a # character is ignored */5 * * * *
//     /usr/local/bin/backup 00 01 * * * /usr/bin/security_upgrades
//
// Running `cronrun` on the example file will run `/usr/local/bin/backup` if the
// current minute is divisible by 5, and will run `/usr/bin/security_upgrades` if
// the time is `01:00`.
//
// Running from cron
//
// Run `cronrun` from your system crontab, as you would any normal cron job.
//
// Example cron job to run `cronrun` on a given file:
//
//     * * * * * /usr/local/bin/cronrun /var/www/mysite/crontab
//
package cronrun
