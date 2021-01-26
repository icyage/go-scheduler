package scheduler

import (
	_ "github.com/lib/pq"
	"github.com/obsurvive/voyager/log"
	"github.com/obsurvive/voyager/utils"
	"os"
	"os/signal"
	"syscall"
)

var (
	signals    chan os.Signal
	jobManager *JobManager
)

// CreateTable creates the table for storing jobs on database.
func CreateTable() error {

	err, db := utils.GetDB()
	if err != nil {
		return err
	}

	defer db.Close()
	t := &Table{DB: db, Name: cfg.GetString("scheduler_table_name")}
	return t.Create()
}

func Run() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	err, db := utils.GetDB()
	if err != nil {
		log.Error(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	table := &Table{DB: db, Name: cfg.GetString("scheduler_table_name")}

	scheduler := NewScheduler(table, cfg.GetFloat64("randomize_factor"))

	go scheduler.Run()
	defer func() {
		scheduler.Stop()
		<-scheduler.NotifyDone()
	}()

	signalReceived := <-signals
	log.WithFields(log.Fields{
		"EventName": "server_start",
		"Signal":    signalReceived,
	}).Infof("Shutting down server: %s", signalReceived)
}

func GetJobManager() (*JobManager) {

	if jobManager == nil {
		err, db := utils.GetDB()
		if err != nil {
			//ToDo: log something
		}
		t := &Table{DB: db, Name: cfg.GetString("scheduler_table_name")}
		s := NewScheduler(t, cfg.GetFloat64("randomize_factor"))
		jobManager = NewJobManager(t, s)
		return jobManager
	} else {
		return jobManager
	}
}
