package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Init("debug", []string{"stdout"}) //[]string{"stdout", "test.log"}

	Info("Test log.Info", " value is ", 10)
	Infof("Test log.Infof %d", 10)
	Infow("Test log.Infow", "value", 10)
	Debugf("Test log.Debugf %d", 10)
	Error("Test log.Error", " value is ", 10)
	Errorf("Test log.Errorf %d", 10)
	Errorw("Test log.Errorw", "value", 10)
	Warnf("Test log.Warnf %d", 10)
	Warnw("Test log.Warnw", "value", 10)
}
