package messaging

import (
	"github.com/streadway/amqp"
)

type Connection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func GetMessagingConnection(url string) (Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return Connection{}, err
	}

	ch, err := conn.Channel()
	return Connection{
		Connection: conn,
		Channel:    ch,
	}, err
}

func (r Connection) Close() {
	r.Connection.Close()
	r.Channel.Close()
}

func (r *Connection) PublishQueue(data []byte, queueName string) error {

	queue, err := r.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return err
	}

	err = r.Channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})

	if err != nil {
		return err
	}

	return nil
}

func (r *Connection) ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {

	queue, err := r.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return nil, err
	}

	msgs, err := r.Channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		return nil, err
	}

	return msgs, nil
}
