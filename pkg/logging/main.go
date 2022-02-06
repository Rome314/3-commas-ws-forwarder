package logging

import (
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

type Entry struct {
	*log.Entry
}

func SetDebug(debug bool) {
	if debug {
		logger.SetLevel(log.DebugLevel)

	} else {
		logger.SetLevel(log.InfoLevel)
	}
}

func (e *Entry) WithMethod(method string) *Entry {
	return &Entry{e.Entry.WithField("method", method)}
}
func (e *Entry) WithPlace(place string) *Entry {
	return &Entry{e.Entry.WithField("place", place)}
}

func GetLogger(module string) *Entry {
	return &Entry{Entry: logger.WithFields(log.Fields{
		"module": module,
	})}
}

func init() {

	logger = log.New()
	logger.SetFormatter(&log.TextFormatter{})

}
