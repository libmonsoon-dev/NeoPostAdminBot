package logrus

import (
	log "github.com/sirupsen/logrus"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
)

type factory struct {
	l *log.Logger
}

func (f factory) New(componentName string) logger.Logger {
	return entity{f.l.WithField("component", componentName)}
}

func NewFactory() logger.Factory {
	l := log.New()
	l.SetLevel(log.TraceLevel)
	l.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	return factory{l}
}
