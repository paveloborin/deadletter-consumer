package consumer

import (
	"fmt"

	"github.com/paveloborin/deadletter-consumer/pkg/flags"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
	consume func(amqp.Delivery)
}

func NewConsumer(config *flags.Config, consume func(amqp.Delivery)) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		done:    make(chan error),
		consume: consume,
	}

	var err error

	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.AmqpUser, config.AmqpPassword, config.AmqpHost, config.AmqpPort)
	if c.conn, err = amqp.Dial(amqpURI); err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		config.DeadLetterExchangeName,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %s", err)
	}

	dlQueue, err := c.channel.QueueDeclare(
		config.DeadLetterQueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("queue Declare: %s", err)
	}

	if err = c.channel.QueueBind(
		dlQueue.Name,
		"",
		config.DeadLetterExchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		config.ExchangeName,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %s", err)
	}

	queue, err := c.channel.QueueDeclare(
		config.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-dead-letter-exchange": config.DeadLetterExchangeName},
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare: %s", err)
	}

	if err = c.channel.QueueBind(
		queue.Name,
		config.RoutingKey,
		config.ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %s", err)
	}

	deliveries, err := c.channel.Consume(
		queue.Name,
		c.tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %s", err)
	}

	go handle(deliveries, c.done, c.consume)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error, consume func(amqp.Delivery)) {
	for d := range deliveries {
		consume(d)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
