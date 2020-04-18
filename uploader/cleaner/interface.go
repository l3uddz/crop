package cleaner

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	FindHidden(*config.UploaderHidden, *logrus.Entry) ([]pathutils.Path, []pathutils.Path, error)
}
