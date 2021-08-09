package rabbitmq

import (
	"Cube-back/log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error) {
	if err != nil {
		log.Error(err)
	}
}

func init() {
	conn, err := amqp.Dial("amqp://guest:guest@1.15.111.85:15672/admin")
	log.Error(err)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"cube",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err)

	q, err := ch.QueueDeclare(
		"blog", // name
		false,  // durable
		false,  // delete when unused
		true,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err)

	if len(os.Args) < 2 {
		log.Info(os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		err = ch.QueueBind(
			q.Name, // queue name
			s,      // routing key
			"cube", // exchange
			false,
			nil)
		log.Error(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	log.Error(err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
