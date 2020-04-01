package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func init() {
	if Log == nil {
		Log = logrus.New()
	}
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
}
