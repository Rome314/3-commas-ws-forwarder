package api

type listener struct {
	conn    *wsConn
	output  chan Message
	closeCh chan struct{}
}

func (l *listener) listen() {
	for {
		select {
		case _ = <-l.closeCh:
			return
		default:
			msg := Message{}
			if err := l.conn.ReadJson(&msg); err != nil {
				continue
			}
			l.output <- msg
		}

	}
}
func (l *listener) close() {
	l.closeCh <- struct{}{}
	close(l.output)
}
