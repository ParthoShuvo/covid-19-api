package log4u

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var logFile *os.File

type customLogFormatter struct {
	log.TextFormatter
}

func (cf *customLogFormatter) Format(entry *log.Entry) ([]byte, error) {
	function, file := cf.CallerPrettyfier(entry.Caller)
	return []byte(fmt.Sprintf("[%s] [%s] [%s:%s] => %s\n",
		entry.Level.String(),
		entry.Time.Format(cf.TimestampFormat),
		file,
		function,
		entry.Message,
	)), nil
}

const defaultLevel = "DEBUG"

func logLevel(level string) log.Level {
	for _, logLevel := range log.AllLevels {
		if logLevel.String() == level {
			return logLevel
		}
	}
	return log.New().Level
}

// ConfigureLogging configures log properties and file
func ConfigureLogging(filename string, level string) {
	if filename == "" {
		log.SetLevel(logLevel(defaultLevel))
	} else {
		var err error
		if logFile, err = os.Create(filename); err != nil {
			log.Fatalf("failed to create file %s: %v", filename, err)
			return
		}
		log.SetLevel(logLevel(level))
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}
	log.SetReportCaller(true)
	log.SetFormatter(formatter())
}

func formatter() log.Formatter {
	return &customLogFormatter{
		log.TextFormatter{
			ForceColors:      false,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				_, filename := path.Split(f.File)
				return fmt.Sprintf("%s:%d", f.Function, f.Line), filename
			},
		},
	}
}

func ContainsLogDebug(level string) bool {
	return level == "DEBUG"
}

// CloseLog closing the log file
func CloseLog() {
	if logFile != nil {
		logFile.Close()
	}
}
