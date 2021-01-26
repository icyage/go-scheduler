package scheduler

import (
	"encoding/json"
	"fmt"
	"time"
)

// Job is the record stored in jobs table.
// Primary key for the table is JobKey.
type Job struct {
	JobKey
	// Interval is the duration between each job.
	Interval time.Duration
	// NextRun is the next run time of the job, stored in UTC.
	NextRun time.Time
}

type JobKey struct {
	// ID of the probe to run
	ID string
}

// String returns the job in human-readable form.
func (j *Job) String() string {
	return fmt.Sprintf("Job{%q, %q, %s}", j.ID, j.Interval, j.NextRun.String()[:23])
}

// Remaining returns the remaining time to the job's next run time.
func (j *Job) Remaining() time.Duration {
	return j.NextRun.Sub(time.Now().UTC())
}

func (j *Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID       string        `json:"id"`
		Interval time.Duration `json:"interval"`
		NextRun  string        `json:"next_run"`
	}{
		ID:       j.ID,
		Interval: j.Interval / time.Second,
		NextRun:  j.NextRun.Format(time.RFC3339),
	})
}
