package rabbitmq

import (
	"Cube-back/log"
	"Cube-back/models/common/configure"
	"github.com/streadway/amqp"
)

var MessageQueue *Mq

type Mq struct {
	RabbitmqIp       string
	RabbitmqPort     string
	RabbitmqPassword string
	RabbitmqUser     string
	conn             *amqp.Connection
}

func (m *Mq) MessageSend(queueName, message string) {
	m.channelCreate(queueName, message)
}

func (m *Mq) channelCreate(queueName, message string) {
	ch, err := m.conn.Channel()
	if err != nil {
		log.Error(err)
		return
	}
	exchangeDeclare(ch, queueName, message)
}

func exchangeDeclare(ch *amqp.Channel, queueName, message string) {
	err := ch.ExchangeDeclare(
		"cube",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Error(err)
		return
	}
	queueDeclare(ch, queueName, message)
}

func queueDeclare(ch *amqp.Channel, queueName, message string) {
	_, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Error(err)
		return
	}
	queueBuild(ch, queueName, message)
}

func queueBuild(ch *amqp.Channel, queueName, message string) {
	err := ch.QueueBind(
		queueName, // queue name
		queueName, // routing key
		"cube",    // exchange
		false,
		nil)
	if err != nil {
		log.Error(err)
		return
	}
	messagePublish(ch, queueName, message)
}

func messagePublish(ch *amqp.Channel, queueName, message string) {
	err := ch.Publish(
		"cube",    // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Error(err)
		ch.Close()
		return
	}
	ch.Close()

}
func (m *Mq) mqStart() {
	mq := new(Mq)
	configure.Get(&mq)
	url := "amqp://" + mq.RabbitmqUser + ":" + mq.RabbitmqPassword + "@" + mq.RabbitmqIp + "/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Error(err)
		return
	}

	m.conn = conn
}

func init() {
	MessageQueue = new(Mq)
	MessageQueue.mqStart()
}
