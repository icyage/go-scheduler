package scheduler

import (
	"database/sql"
	"github.com/obsurvive/voyager/log"
	"github.com/obsurvive/voyager/utils"
	"github.com/streadway/amqp"
	"math/rand"
	"sync"
	"time"
)

type scheduler struct {
	table               *Table
	baseURL             string
	randomizationFactor float64
	// to stop scheduler goroutine
	stop chan struct{}
	// will be closed when scheduler goroutine is stopped
	stopped chan struct{}
	// to wake up scheduler when a new job is scheduled or cancelled
	wakeUp      chan struct{}
	runningJobs map[JobKey]struct{}
	m           sync.Mutex
	wg          sync.WaitGroup
}

func NewScheduler(t *Table, randomizationFactor float64) *scheduler {
	s := &scheduler{
		table:               t,
		randomizationFactor: randomizationFactor,
		stop:                make(chan struct{}),
		stopped:             make(chan struct{}),
		wakeUp:              make(chan struct{}, 1),
		runningJobs:         make(map[JobKey]struct{}),
	}
	return s
}

func (s *scheduler) WakeUp(debugMessage string) {
	select {
	case s.wakeUp <- struct{}{}:
		log.WithFields(log.Fields{
			"EventName": "scheduler_waku_up",
			"DebugMessage":    debugMessage,
		}).Debugf("notifying scheduler:", debugMessage)
	default:
	}
}

func (s *scheduler) NotifyDone() <-chan struct{} {
	return s.stopped
}

func (s *scheduler) Stop() {
	close(s.stop)
}

func (s *scheduler) Running() int {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.runningJobs)
}

// Run runs a loop that reads the next Job from the queue and executees it in it's own goroutine.
func (s *scheduler) Run() {

	err, queue := utils.GetQueue()
	if err != nil {
		log.Fatal(err)
	}

	defer queue.Close()
	defer close(s.stopped)

	channel, err := queue.Channel()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", "direct", "voyager-probe")
	if err := channel.ExchangeDeclare(
		cfg.GetString("amqp_exchange_name"),     // name
		"direct", // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		log.Fatal(err)
	}

	_, err = channel.QueueDeclare(cfg.GetString("amqp_queue_name"), true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err =channel.QueueBind(cfg.GetString("amqp_queue_name"), cfg.GetString("amqp_key_name"), cfg.GetString("amqp_exchange_name"), false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		var after <-chan time.Time

		job, err := s.table.Front()
		if err != nil {
			if err == sql.ErrNoRows {
				log.WithFields(log.Fields{
					"EventName": "scheduler_run",
					"DebugMessage":    "NoJobs",
				}).Debug("no scheduled jobs in the table")
			} else {
				log.WithFields(log.Fields{
					"EventName": "scheduler_run_fatal",
					"Error": err,
				}).Fatal("DB Error")
				time.Sleep(time.Second)
				continue
			}
		} else {
			remaining := job.Remaining()
			after = time.After(remaining)
			log.WithFields(log.Fields{
				"EventName": "scheduler_run",
				"DebugMessage":    "got_job",
				"Job": job,
				"Remaining": remaining,
			}).Debugf("next job:", job, "remaining:", remaining)
		}

		// Sleep until the next job's run time or the webserver's wakes us up.
		select {
		case <-after:
			log.WithFields(log.Fields{
				"EventName": "scheduler_run",
				"DebugMessage":    "sleep_finished",
			}).Debug("job sleep time finished")


			log.Printf("declared Exchange, publishing %dB body (%q)", len(job.ID), job.ID)
			if err = channel.Publish(
				cfg.GetString("amqp_exchange_name"),   // publish to an exchange
				cfg.GetString("amqp_key_name"), // routing to 0 or more queues
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					Headers:         amqp.Table{},
					ContentType:     "text/plain",
					ContentEncoding: "",
					Body:            []byte(job.ID),
					DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
					Priority:        0,              // 0-9
				},
			); err != nil {
				log.Error(err)
			}

			add := job.Interval
			if s.randomizationFactor > 0 {
				// Add some randomization to periodic tasks.
				add = randomize(add, s.randomizationFactor)
			}

			job.NextRun = time.Now().UTC().Add(add)

			err = s.table.UpdateNextRun(job)

			if err != nil {
				log.Error(err)
			}

			log.Debugf("scheduling job", job)
			/*if err = s.execute(job); err != nil {
				log.WithFields(log.Fields{
					"EventName": "scheduler_run_error",
					"Error": err,
				}).Errorf("Error queueing job", err)
				time.Sleep(time.Second)
			}*/
		case <-s.wakeUp:
			log.WithFields(log.Fields{
				"EventName": "scheduler_run",
				"DebugMessage":    "woke_up",
			}).Debug("woken up from sleep by notification")
			continue
		case <-s.stop:
			log.WithFields(log.Fields{
				"EventName": "scheduler_run",
				"DebugMessage":    "quit",
			}).Debug("came quit message")
			s.wg.Wait()
			return
		}
	}
}

func randomize(d time.Duration, f float64) time.Duration {
	delta := time.Duration(f * float64(d))
	return d - delta + time.Duration(float64(2*delta)*rand.Float64())
}