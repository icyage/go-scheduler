package utils

import (
	"github.com/streadway/amqp"
)

var (
	queue *amqp.Connection
)

func GetQueue() (error, *amqp.Connection) {

	if queue == nil {
		queue, err := amqp.Dial(cfg.GetString("amqp_dsn"))
		if err != nil {
			return err, nil
		} else {
			return nil, queue
		}
	}
	return nil, queue
}
