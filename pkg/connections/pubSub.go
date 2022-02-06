package connections

import (
	"context"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-redis/redis/v8"
	redisPubSub "github.com/minghsu0107/watermill-redistream/pkg/redis"
	"github.com/rome314/3-commas-ws-forwarder/internal/sender"
)

type RedisPubSub struct {
	Pub message.Publisher
	Sub message.Subscriber
}

type RedisPubSubConfg struct {
}

func GetRedisPubSub(ctx context.Context, client redis.UniversalClient, consumerGroup string) (pubSub *RedisPubSub, err error) {
	pubSubMarshaler := sender.RedisMarshaller{}
	sub, err := redisPubSub.NewSubscriber(
		ctx,
		redisPubSub.SubscriberConfig{
			Consumer:      consumerGroup,
			ConsumerGroup: consumerGroup,
		},
		client,
		pubSubMarshaler,
		nil,
	)
	if err != nil {
		err = errors.WithMessage(err, "creating sub")
		return
	}

	pub, err := redisPubSub.NewPublisher(
		ctx,
		redisPubSub.PublisherConfig{},
		client,
		pubSubMarshaler,
		nil,
	)
	if err != nil {
		err = errors.WithMessage(err, "creating pub")
		return
	}
	return &RedisPubSub{
		Pub: pub,
		Sub: sub,
	}, nil

}
