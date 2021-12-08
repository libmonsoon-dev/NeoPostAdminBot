package logrus

import (
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
	"github.com/sirupsen/logrus"
)

var _ logger.Logger = (*logrus.Logger)(nil)
