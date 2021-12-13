package logrus

import (
	log "github.com/sirupsen/logrus"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
)

type entity struct {
	*log.Entry
}

func (e entity) IsLevelEnabled(level logger.Level) bool {
	return e.Logger.IsLevelEnabled(log.Level(level))
}
