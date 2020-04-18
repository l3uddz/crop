package checker

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	Check(*config.UploaderCheck, *logrus.Entry, []pathutils.Path, uint64) (bool, error)
	CheckFile(*config.UploaderCheck, *logrus.Entry, pathutils.Path, uint64) (bool, error)
}
