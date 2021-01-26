package scheduler

import (
	"errors"
	"github.com/obsurvive/voyager/log"
	"time"
)

type JobManager struct {
	table     *Table
	scheduler *scheduler
}

var errInvalidArgs = errors.New("invalid arguments")

func NewJobManager(t *Table, s *scheduler) *JobManager {
	return &JobManager{
		table:     t,
		scheduler: s,
	}
}

// Get returns the job with path and body.
func (m *JobManager) Get(ID string) (*Job, error) {
	return m.table.Get(ID)
}

// Schedule inserts a new job to the table or replaces existing one.
// Returns the created or replaced job.
func (m *JobManager) Schedule(ID string, interval *time.Duration) (*Job, error) {
	job := &Job{JobKey: JobKey{
		ID: ID,
	}}

	if interval == nil {
		return nil, errInvalidArgs
	}
	job.Interval = *interval
	job.NextRun = time.Now().UTC().Add(*interval)

	err := m.table.Insert(job)
	if err != nil {
		return nil, err
	}
	m.scheduler.WakeUp("new job")
	log.WithFields(log.Fields{
		"EventName":    "job_manager",
		"DebugMessage": "job_scheduled",
		"Job":          job,
	}).Debugf("job is scheduled:", job)
	return job, nil
}

// Trigger runs the job immediately and resets it's next run time.
func (m *JobManager) Trigger(ID string) (*Job, error) {
	job, err := m.table.Get(ID)
	if err != nil {
		return nil, err
	}
	job.NextRun = time.Now().UTC()
	if err := m.table.Insert(job); err != nil {
		return nil, err
	}
	m.scheduler.WakeUp("job is triggered")
	log.WithFields(log.Fields{
		"EventName":    "job_manager",
		"DebugMessage": "job_triggered",
		"Job":          job,
	}).Debugf("job is triggered:", job)
	return job, nil
}

// Cancel deletes the job with path and body.
func (m *JobManager) Cancel(ID string) error {
	err := m.table.Delete(ID)
	if err != nil {
		return err
	}
	m.scheduler.WakeUp("job cancelled")
	log.WithFields(log.Fields{
		"EventName":    "job_manager",
		"DebugMessage": "job_canceled",
		"Job":          ID,
	}).Debugf("job is canceled:", ID)
	return nil
}

// Running returns the number of running jobs currently.
func (m *JobManager) Running() int {
	return m.scheduler.Running()
}

// Total returns the count of all jobs in jobs table.
func (m *JobManager) Total() (int64, error) {
	return m.table.Count()
}
