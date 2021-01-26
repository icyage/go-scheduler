package scheduler

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/obsurvive/voyager/log"
	"time"
)

type Table struct {
	DB   *sql.DB
	Name string
}

const createTableSQL = "CREATE TABLE %s (id VARCHAR(255) NOT NULL, interval INT CHECK (interval > 0) NOT NULL, next_run TIMESTAMP(0) NOT NULL, PRIMARY KEY (id))"
const createTableIndex = "CREATE INDEX idx_next_run ON %s (next_run)"

var (
	ErrExist    = errors.New("job already exists")
	ErrNotExist = errors.New("job does not exist")
)

// Create jobs table.
func (t *Table) Create() error {
	sqlQuery := fmt.Sprintf(createTableSQL, t.Name)
	_, err := t.DB.Exec(sqlQuery)

	if err != nil {
		return err
	}

	sqlQuery = fmt.Sprintf(createTableIndex, t.Name)
	_, err = t.DB.Exec(sqlQuery)

	return err
}

// Get returns a job with body and path from the table.
func (t *Table) Get(ID string) (*Job, error) {
	var j Job
	var interval uint32

	err := t.DB.QueryRow("SELECT id, interval, next_run FROM "+t.Name+" WHERE id = $1", ID).Scan(&j.ID, &interval, &j.NextRun)
	log.Debug(err)
	if err != nil {
		return nil, err
	}

	j.Interval = time.Duration(interval) * time.Second
	return &j, nil
}

// Insert the job to to scheduler table.
func (t *Table) Insert(j *Job) error {
	_, err := t.DB.Exec("INSERT INTO "+t.Name+"(id, interval, next_run) VALUES ($1, $2, $3)", j.ID, j.Interval.Seconds(), j.NextRun)
	return err
}

// Delete the job from scheduler table.
func (t *Table) Delete(ID string) error {
	result, err := t.DB.Exec("DELETE FROM "+t.Name+" WHERE id = $1", ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotExist
	}
	return nil
}

// Front returns the next scheduled job from the table.
func (t *Table) Front() (*Job, error) {

	row := t.DB.QueryRow("SELECT id, interval, next_run FROM " + t.Name + " ORDER BY next_run ASC LIMIT 1")
	var j Job
	var interval uint32
	err := row.Scan(&j.ID, &interval, &j.NextRun)
	if err != nil {
		return nil, err
	}
	j.Interval = time.Duration(interval) * time.Second
	return &j, nil
}

// UpdateNextRun sets next_run to now+interval.
func (t *Table) UpdateNextRun(j *Job) error {
	_, err := t.DB.Exec("UPDATE "+t.Name+" SET next_run=$1 WHERE id = $2", j.NextRun, j.ID)
	return err
}

// Count returns the count of scheduled jobs in the table.
func (t *Table) Count() (int64, error) {
	var count int64
	return count, t.DB.QueryRow("SELECT COUNT(*) FROM " + t.Name).Scan(&count)
}
