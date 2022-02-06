package sender

import (
	"net/url"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rome314/3-commas-ws-forwarder/internal/api"
)

type Config struct {
	Topic         string
	Url           string
	ConsumerGroup string
}

func (c Config) Valid() error {
	if c.Topic == "" {
		return errors.New("topic not provided")
	}
	if _, err := url.ParseRequestURI(c.Url); err != nil {
		return errors.WithMessage(err, "invalid url")
	}

	if c.ConsumerGroup == "" {
		return errors.New("consumer group not provided")
	}

	return nil
}

type CreateInput struct {
	Pub    message.Publisher
	Sub    message.Subscriber
	Cfg    Config
	Events <-chan api.Message
}
