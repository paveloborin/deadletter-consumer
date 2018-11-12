package flags

type Config struct {
	PrettyLogging bool `long:"pretty-logging" description:"Pretty logging"`

	AmqpHost     string `long:"amqp-host" env:"AMQP_HOST" required:"true"`
	AmqpPort     int    `long:"amqp-port" env:"AMQP_PORT" required:"true"`
	AmqpUser     string `long:"amqp-user" env:"AMQP_USER" required:"true"`
	AmqpPassword string `long:"amqp-password" env:"AMQP_PASSWORD" required:"true"`

	ExchangeName string `long:"exchange-name" env:"EXCHANGE_NAME" required:"true"`
	QueueName    string `long:"queue-name" env:"QUEUE_NAME" required:"true"`
	RoutingKey   string `long:"routing-key" env:"ROUTING_KEY" required:"false" default:""`

	DeadLetterExchangeName string `long:"dead-letter-exchange-name" env:"DEAD_LETTER_EXCHANGE_NAME" required:"true"`
	DeadLetterQueueName    string `long:"dead-letter-queue-name" env:"DEAD_LETTER_QUEUE_NAME" required:"true"`
}
