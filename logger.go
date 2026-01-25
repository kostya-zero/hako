package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var l *log.Logger

func PrepareLogger() {
	l = log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.RFC1123,
	})
}
