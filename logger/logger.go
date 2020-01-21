package logger

import (
	log "github.com/sirupsen/logrus"
)

var Behaviour *log.Logger

func Init() {
	Behaviour = log.New()

	Behaviour.Formatter = &log.TextFormatter{
		ForceColors:               true,
		FullTimestamp:             true,
		TimestampFormat:           "",
	}
	Behaviour.Level = log.DebugLevel

}