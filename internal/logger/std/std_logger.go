package std

import (
	"os"

	"github.com/number571/go-peer/pkg/logger"
)

const (
	CLogInfo = "info"
	CLogWarn = "warn"
	CLogErro = "erro"
)

func NewStdLogger(pLogging ILogging, pLogFunc logger.ILogFunc) logger.ILogger {
	return logger.NewLogger(
		defaultSettings(pLogging),
		pLogFunc,
	)
}

func defaultSettings(pLogging ILogging) logger.ISettings {
	sett := &logger.SSettings{}
	if pLogging.HasInfo() {
		sett.FInfo = os.Stdout
	}
	if pLogging.HasWarn() {
		sett.FWarn = os.Stdout
	}
	if pLogging.HasErro() {
		sett.FErro = os.Stderr
	}
	return logger.NewSettings(sett)
}
