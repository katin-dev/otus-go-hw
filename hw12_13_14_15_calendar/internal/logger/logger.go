package logger

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	file string
	logg *logrus.Logger
}

func New(file, level, formatter string) (*Logger, error) {
	log := logrus.New()

	switch file {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		file, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			return nil, fmt.Errorf("invalid log filename: %w", err)
		}
	}

	levelID, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	log.SetLevel(levelID)

	switch formatter {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text_simple":
		log.SetFormatter(&SimpleTextFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{})
	}

	logger := &Logger{
		file: file,
		logg: log,
	}

	return logger, nil
}

func (l *Logger) Debug(msg string, params ...interface{}) {
	l.logg.Debugf(msg, params...)
}

func (l *Logger) Info(msg string, params ...interface{}) {
	l.logg.Infof(msg, params...)
}

func (l *Logger) Warn(msg string, params ...interface{}) {
	l.logg.Warnf(msg, params...)
}

func (l *Logger) Error(msg string, params ...interface{}) {
	l.logg.Errorf(msg, params...)
}

func (l *Logger) LogHttpRequest(r *http.Request, code, length int) {
	l.logg.Infof(
		"%s\t%s\t%s\t%s\t%d\t%d\t\"%s\"",
		r.RemoteAddr,
		r.Method,
		r.URL.String(),
		r.Proto,
		code,
		length,
		r.UserAgent(),
	)
}

type SimpleTextFormatter struct{}

func (f *SimpleTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s\t%s\n", entry.Level, entry.Message)

	return []byte(msg), nil
}
