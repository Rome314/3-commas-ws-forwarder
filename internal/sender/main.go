package sender

import (
	"context"
	"encoding/json"
	"net/http"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rome314/3-commas-ws-forwarder/internal/api"
	"github.com/rome314/3-commas-ws-forwarder/pkg/logging"
)

type sndr struct {
	logger *logging.Entry
	pub    message.Publisher
	sub    message.Subscriber
	events <-chan api.Message
	cfg    Config
}

func NewSender(logger *logging.Entry, input CreateInput) *sndr {
	return &sndr{
		logger: logger,
		pub:    input.Pub,
		sub:    input.Sub,
		events: input.Events,
		cfg:    input.Cfg,
	}
}

func (s *sndr) publisher(ctx context.Context, close chan struct{}) {
	logger := s.logger.WithMethod("publisher")
	select {
	case _ = <-ctx.Done():
		close <- struct{}{}
		return
	default:
		for msg := range s.events {
			if msg.Type == api.PingMessageType {
				logger.Debug("ping")
				continue
			}
			logger.Infof("Handled evet: %s", msg.Type)
			bts, err := json.Marshal(msg)
			if err != nil {
				logger.WithPlace("unmarshal_message").Error(err)
				continue
			}
			toSend := message.NewMessage(watermill.NewUUID(), bts)
			if err = s.pub.Publish(s.cfg.Topic, toSend); err != nil {
				logger.WithPlace("publish_message").Error(err)

			}
		}
		close <- struct {
		}{}
	}
}

func (s *sndr) subscriber(msgs <-chan *message.Message) {
	logger := s.logger.WithMethod("subscriber")
	client := &http.Client{}
	for msg := range msgs {
		parsed := api.Message{}
		if e := json.Unmarshal(msg.Payload, &parsed); e != nil {
			logger.WithPlace("unmarshal_message").Error(e)
			msg.Nack()
			continue
		}
		if _, e := client.Post(s.cfg.Url, "application/json", parsed.Reader()); e != nil {
			logger.WithPlace("send_webhook").Error(e)
			msg.Nack()
			continue
		}
		msg.Ack()
	}
}

func (s *sndr) Run(ctx context.Context) (closeChan chan struct{}, err error) {
	msgs, err := s.sub.Subscribe(ctx, s.cfg.Topic)
	if err != nil {
		return nil, errors.WithMessage(err, "subscribing to topic")
	}
	closeChan = make(chan struct{})
	go s.subscriber(msgs)
	go s.publisher(ctx, closeChan)
	return closeChan, nil
}
