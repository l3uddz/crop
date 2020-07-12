package checker

import (
	"github.com/l3uddz/crop/config"
	"github.com/l3uddz/crop/pathutils"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	Check(*config.UploaderCheck, *logrus.Entry, []pathutils.Path, uint64) (*Result, error)
	CheckFile(*config.UploaderCheck, *logrus.Entry, pathutils.Path, uint64) (bool, error)
	RcloneParams(check *config.UploaderCheck, entry *logrus.Entry) []string
}

type Result struct {
	Passed bool
	Info   interface{}
}
