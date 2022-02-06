package api

import (
	"sync"

	"emperror.dev/errors"
	"github.com/gorilla/websocket"
)

type wsConn struct {
	conn *websocket.Conn
	mx   *sync.Mutex
}

func connectToWs() (wc *wsConn, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		err = errors.WithMessage(err, "connecting to ws")
		return
	}

	wc = newWsConn(conn)
	return wc, nil
}

func newWsConn(conn *websocket.Conn) *wsConn {
	return &wsConn{
		conn: conn,
		mx:   &sync.Mutex{},
	}
}

func (w *wsConn) IsClosed() bool {
	w.mx.Lock()
	defer w.mx.Unlock()
	return w.conn.WriteMessage(websocket.PingMessage, []byte{}) != nil

}

func (w *wsConn) ReadJson(value interface{}) error {
	w.mx.Lock()
	defer w.mx.Unlock()
	return w.conn.ReadJSON(value)
}

func (w *wsConn) WriteJson(value interface{}) error {
	w.mx.Lock()
	defer w.mx.Unlock()
	return w.conn.WriteJSON(value)
}
