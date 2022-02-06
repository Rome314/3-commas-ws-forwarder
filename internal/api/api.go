package api

import (
	"context"
	"time"

	"emperror.dev/errors"
)

const wsEndpoint = "wss://ws.3commas.io/websocket"

type Config struct {
	ApiKey    string
	SecretKey string
}

type apiImpl struct {
	cfg       Config
	listeners map[*wsConn]*listener
}

func New(config Config) *apiImpl {
	a := &apiImpl{cfg: config, listeners: map[*wsConn]*listener{}}
	go a.connsChecker()
	return a
}
func (a *apiImpl) connsChecker() {
	for {
		time.Sleep(time.Second * 5)
		for conn, _ := range a.listeners {
			if conn.IsClosed() {
				panic("websocket closed")
			}
		}
	}
}

func (a *apiImpl) SubscribeChannel(ctx context.Context, channel string) (messages <-chan Message, err error) {

	ident, err := getIdentifier(channel, a.cfg)
	if err != nil {
		err = errors.WithMessage(err, "getting identifier")
		return
	}

	wc, err := connectToWs()
	if err != nil {
		return
	}

	command := Command{
		Command:    subscribeCommand,
		Identifier: ident.String(),
	}

	if err = wc.WriteJson(command); err != nil {
		err = errors.WithMessage(err, "sending cub message")
		return
	}

	loop := true
	for loop {
		select {
		case _ = <-ctx.Done():
			err = errors.New("subscribe timeout")
			return
		default:
			msg := Message{}
			if err = wc.ReadJson(&msg); err != nil {
				err = errors.WithMessage(err, "reading message")
				return
			}
			if msg.Type == ConfirmSubscriptionMessageType {
				loop = false
			}
		}
	}

	msgs := make(chan Message)
	l := &listener{
		conn:    wc,
		output:  msgs,
		closeCh: make(chan struct{}),
	}
	a.listeners[wc] = l
	go l.listen()

	return msgs, nil
}
