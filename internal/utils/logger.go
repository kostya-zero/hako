package utils

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var L *log.Logger

func PrepareLogger() {
	L = log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.RFC1123,
	})
}
