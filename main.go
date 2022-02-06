package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"emperror.dev/emperror"
	"github.com/rome314/3-commas-ws-forwarder/internal/api"
	"github.com/rome314/3-commas-ws-forwarder/internal/config"
	"github.com/rome314/3-commas-ws-forwarder/internal/sender"
	"github.com/rome314/3-commas-ws-forwarder/pkg/connections"
	"github.com/rome314/3-commas-ws-forwarder/pkg/logging"
)

func main() {
	logger := logging.GetLogger("main")

	logger.Info("Preparing config...")
	cfg := config.GetConfig()
	emperror.Panic(cfg.Valid())
	logging.SetDebug(cfg.Debug)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	logger.Info("Preparing DB connections...")

	logger.Info("Preparing Redis connection...")
	redisConn, err := connections.GetRedisConnection(ctx, connections.RedisConfig{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
		Db:       cfg.Redis.Db,
	})
	emperror.Panic(err)

	logger.Info("Configuring internal modules...")
	pubSub, err := connections.GetRedisPubSub(ctx, redisConn.Connection, cfg.App.ConsumerGroup)
	emperror.Panic(err)

	commasApi := api.New(api.Config{
		ApiKey:    cfg.Api.Key,
		SecretKey: cfg.Api.Secret,
	})

	timeoutCtx, _ := context.WithTimeout(ctx, time.Second*5)
	apiEvents, err := commasApi.SubscribeChannel(timeoutCtx, api.DealsChannel)
	emperror.Panic(err)

	senderInput := sender.CreateInput{
		Pub:    pubSub.Pub,
		Sub:    pubSub.Sub,
		Cfg:    cfg.App,
		Events: apiEvents,
	}

	sndr := sender.NewSender(logging.GetLogger("sender"), senderInput)

	logger.Info("Running service...")
	senderStop, err := sndr.Run(ctx)
	emperror.Panic(err)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Kill, os.Interrupt)
	go func() {
		<-senderStop
		stop <- os.Kill
	}()

	<-stop
}
