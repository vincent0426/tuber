package mq

import (
	"log"

	"github.com/streadway/amqp"
)

type rabbitMQ struct {
	amqpConnection    *amqp.Connection
	rabbitMQChannel   *RabbitMQChannelImpl
	delayQueueName    string
	delayExchangeName string
	delayRoutingKey   string
}

type RabbitMQChannel interface {
	declareQueue(qName string) error
	declareDelayExchange(exName string) error
	bindQueue(string, string, string) error
	consume(qName string) (<-chan amqp.Delivery, error)
	SendMsg(string, string, []byte, int) error
	NewMsgReceiver(string, string, string) <-chan amqp.Delivery

	Close() error
}

type RabbitMQChannelImpl struct {
	*amqp.Channel
}

var rabbitmq *rabbitMQ

type Config struct {
	Host              string
	DelayQueueName    string
	DelayExchangeName string
	DelayRoutingKey   string
}

func ConnectToRabbitMQ(cfg Config) (*rabbitMQ, error) {
	conn, err := amqp.Dial(cfg.Host)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	rabbitmq = &rabbitMQ{
		amqpConnection:    conn,
		rabbitMQChannel:   &RabbitMQChannelImpl{ch},
		delayQueueName:    cfg.DelayQueueName,
		delayExchangeName: cfg.DelayExchangeName,
		delayRoutingKey:   cfg.DelayRoutingKey,
	}
	return rabbitmq, err
}

// func NewChannel() {
// 	ch, err := rabbitmq.amqpConnection.Channel()
// 	if err != nil {
// 		log.Fatalf("failed to open a channel: %v", err)
// 	}

// 	rabbitmq.rabbitMQChannel = &RabbitMQChannelImpl{ch}
// 	return
// }

func Close() {
	err := rabbitmq.rabbitMQChannel.Close()
	if err != nil {
		log.Printf("failed to close channel: %v", err)
	}
	err = rabbitmq.amqpConnection.Close()
	if err != nil {
		log.Printf("failed to close connection: %v", err)
	}
}

func (c *RabbitMQChannelImpl) declareQueue(qName string) error {
	_, err := c.QueueDeclare(
		qName, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (c *RabbitMQChannelImpl) declareDelayExchange(exName string) error {
	args := amqp.Table{
		"x-delayed-type": "direct",
	}
	return c.ExchangeDeclare(
		exName,              // name
		"x-delayed-message", // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		args,                // arguments
	)
}

func (c *RabbitMQChannelImpl) bindQueue(qName, exName, routingKey string) error {
	return c.QueueBind(
		qName,      // queue name
		routingKey, // routing key
		exName,     // exchange
		false,      // no-wait
		nil,
	)
}

func (c *RabbitMQChannelImpl) consume(qName string) (<-chan amqp.Delivery, error) {
	return c.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func SendDelayMsg(msg []byte, delayPeriod int) error {
	return rabbitmq.rabbitMQChannel.Publish(
		rabbitmq.delayExchangeName, // exchange
		rabbitmq.delayRoutingKey,   // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
			Headers: amqp.Table{
				"x-delay": delayPeriod,
			},
		})
}

func NewDelayMsgsReceiver() <-chan amqp.Delivery {
	if err := rabbitmq.rabbitMQChannel.declareQueue(rabbitmq.delayQueueName); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	if err := rabbitmq.rabbitMQChannel.declareDelayExchange(rabbitmq.delayExchangeName); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	if err := rabbitmq.rabbitMQChannel.bindQueue(rabbitmq.delayQueueName, rabbitmq.delayExchangeName, rabbitmq.delayRoutingKey); err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	msgReceiver, err := rabbitmq.rabbitMQChannel.consume(rabbitmq.delayQueueName)
	if err != nil {
		log.Fatalf("failed to consume: %v", err)
	}

	return msgReceiver
}
