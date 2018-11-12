package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	pkgFlags "github.com/paveloborin/deadletter-consumer/pkg/flags"
	"github.com/paveloborin/deadletter-consumer/pkg/services/consumer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

var (
	config *pkgFlags.Config
	parser *flags.Parser
)

type Message struct {
	Data map[string]string `json:"data"`
}

func init() {
	zerolog.MessageFieldName = "MESSAGE"
	zerolog.LevelFieldName = "LEVEL"
	log.Logger = log.Output(os.Stderr).With().Str("PROGRAM", "consumer").Logger()

	config = &pkgFlags.Config{}

	parser = flags.NewParser(config, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if config.PrettyLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

}

func main() {
	c, err := consumer.NewConsumer(
		config,
		consumeFunc,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("consumer start failed")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sc

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("error during shutdown")
	}
}

func consumeFunc(message amqp.Delivery) {
	var msg *Message

	err := json.Unmarshal(message.Body, &msg)
	if err != nil {
		log.Warn().Interface("body", message.Body).Msg("Bad message format")
		if err := message.Nack(false, false); err != nil {
			log.Warn().Err(err).Msg("Message nack error")
		}
		return
	}

	//TODO
	log.Info().Interface("mes", msg).Msg("Body")
	if err := message.Ack(false); err != nil {
		log.Warn().Err(err).Msg("Message ack error")
	}
}
