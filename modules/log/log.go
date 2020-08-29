package log

import (
	"github.com/sirupsen/logrus"
)

// Log is a public logger of whole project
var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	//Log.SetReportCaller(true)
}
