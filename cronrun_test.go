package cronrun

import (
	"log"
	"testing"
	"time"
)

func TestNewJob(t *testing.T) {
	cases := []struct {
		input       string
		due         string
		command     string
		errExpected bool
	}{
		{
			input:       "* * * * * foo",
			due:         "* * * * *",
			command:     "foo",
			errExpected: false,
		},
		{
			input:       " * * * * 1 foo bar baz",
			due:         "* * * * 1",
			command:     "foo bar baz",
			errExpected: false,
		},
		{
			input:       "* * * * *",
			due:         "",
			command:     "",
			errExpected: true,
		},
	}
	for _, tc := range cases {
		got, err := NewJob(tc.input)
		if err != nil && !tc.errExpected {
			t.Errorf("NewJob(%q) errored unexpectedly: %v", tc.input, err)
		}
		if err == nil && tc.errExpected {
			t.Errorf("NewJob(%q) did not error as expected", tc.input)
		}
		if got.Due != tc.due {
			t.Errorf("NewJob(%q) => due: %q, want %q", tc.input, got.Due, tc.due)
		}
		if got.Command != tc.command {
			t.Errorf("NewJob(%q) => command: %q, want %q", tc.input, got.Command, tc.command)
		}
	}
}

func TestDueAt(t *testing.T) {
	cases := []struct {
		input string
		now   time.Time
		want  bool
	}{
		{"* * * * *", mustParseTime("2006-01-02T15:04:05Z"), true},
		{"59 * * * *", mustParseTime("2006-01-02T15:04:05Z"), false},
		{"5 * * * *", mustParseTime("2006-01-02T15:04:05Z"), false},
		{"4 * * * *", mustParseTime("2006-01-02T15:04:05Z"), true},
		{"09 12 * * *", mustParseTime("2018-06-02T12:09:05Z"), true},
		{"09 12 * * *", mustParseTime("2018-06-02T12:09:59Z"), true},
		{"09 12 * * *", mustParseTime("2018-06-02T12:08:59Z"), false},
		{"* * * * Mon", mustParseTime("2018-05-08T12:08:59Z"), false},
		{"* * * * Tue", mustParseTime("2018-05-08T12:08:59Z"), true},
		{"* * 1 * *", mustParseTime("2018-05-01T12:08:59Z"), true},
		{"* * 1 * *", mustParseTime("2018-05-02T12:08:59Z"), false},
		{"* * * * Mon-Wed", mustParseTime("2018-05-08T12:08:59Z"), true},
		{"* * * * Thu-Fri", mustParseTime("2018-05-08T12:08:59Z"), false},
		{"* * * 5 *", mustParseTime("2018-05-08T12:08:59Z"), true},
		{"* * * 5 *", mustParseTime("2018-06-08T12:08:59Z"), false},
	}
	for _, tc := range cases {
		j := &Job{tc.input, ""}
		got, err := j.DueAt(tc.now)
		if err != nil {
			t.Errorf("DueAt(%q) at %s errored: %v", tc.input, tc.now.Format(time.RFC3339), err)
		}
		if got != tc.want {
			t.Errorf("DueAt(%q) at %s => %t, want %t", tc.input, tc.now.Format(time.RFC3339), got, tc.want)
		}
	}
	j := &Job{"*bogus*", ""}
	_, err := j.DueAt(time.Now())
	if err == nil {
		t.Errorf("DueAt(bogus data) did not error as expected")
	}
}

func TestRun(t *testing.T) {
	cases := []struct {
		cmd            string
		errExpected    bool
		outputExpected bool
	}{
		{"/bin/bash --version", false, false},
		{"/bin/bash --bogus", true, true},
		{"/bin/bogus --bash", true, false},
	}
	for _, tc := range cases {
		j := &Job{"", tc.cmd}
		output, err := j.Run()
		if err == nil && tc.errExpected {
			t.Errorf("Run(%s) did not error as expected", tc.cmd)
		}
		if err != nil && !tc.errExpected {
			t.Errorf("Run(%s) errored unexpectedly: %v", tc.cmd, err)
		}
		if !tc.outputExpected && len(output) != 0 {
			t.Errorf("Run(%s) wanted no output, got %q", tc.cmd, output)
		}
		if tc.outputExpected && len(output) == 0 {
			t.Errorf("Run(%s) wanted output, got none", tc.cmd)
		}
	}
}

func mustParseTime(ts string) time.Time {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
