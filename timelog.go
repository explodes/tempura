package tempura

import (
	"fmt"
	"log"
	"time"
)

type DurationLogger struct {
	startTime   time.Time
	msg         string
	startLogged bool
}

func LogStart(msg string, args ...interface{}) *DurationLogger {
	logger := &DurationLogger{
		msg: fmt.Sprintf(msg, args...),
	}
	logger.Start()
	return logger
}

func LogDuration(msg string, args ...interface{}) *DurationLogger {
	logger := &DurationLogger{
		startTime: time.Now(),
		msg:       fmt.Sprintf(msg, args...),
	}
	return logger
}

func (l *DurationLogger) Start() {
	l.startTime = time.Now()
	l.startLogged = true
	log.Printf("--> %s", l.msg)
}

func (l *DurationLogger) End() {
	if l == nil {
		return
	}
	if l.startLogged {
		log.Printf("<-- %s (%v)", l.msg, time.Now().Sub(l.startTime))
	} else {
		log.Printf("%s: %v", l.msg, time.Now().Sub(l.startTime))
	}
}
