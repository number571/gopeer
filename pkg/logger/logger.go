package logger

import (
	"fmt"
	"log"
)

var (
	_ ILogger = &sLogger{}
)

type sLogger struct {
	fLogFunc ILogFunc
	fOutInfo *log.Logger
	fOutWarn *log.Logger
	fOutErro *log.Logger
}

const (
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

func NewLogger(pSett ISettings, pLogFunc ILogFunc) ILogger {
	logger := &sLogger{fLogFunc: pLogFunc}

	outInfo := pSett.GetInfoWriter()
	if outInfo != nil {
		logger.fOutInfo = log.New(outInfo, fmt.Sprintf("%s[INFO] %s", colorCyan, colorReset), log.LstdFlags)
	}

	outWarn := pSett.GetWarnWriter()
	if outWarn != nil {
		logger.fOutWarn = log.New(outWarn, fmt.Sprintf("%s[WARN] %s", colorYellow, colorReset), log.LstdFlags)
	}

	outErro := pSett.GetErroWriter()
	if outErro != nil {
		logger.fOutErro = log.New(outErro, fmt.Sprintf("%s[ERRO] %s", colorRed, colorReset), log.LstdFlags)
	}

	return logger
}

func (p *sLogger) PushInfo(pMsg ILogArg) {
	if p.fOutInfo == nil {
		return
	}
	log := p.fLogFunc(pMsg)
	if log == "" {
		return
	}
	p.fOutInfo.Println(log)
}

func (p *sLogger) PushWarn(pMsg ILogArg) {
	if p.fOutWarn == nil {
		return
	}
	log := p.fLogFunc(pMsg)
	if log == "" {
		return
	}
	p.fOutWarn.Println(log)
}

func (p *sLogger) PushErro(pMsg ILogArg) {
	if p.fOutErro == nil {
		return
	}
	log := p.fLogFunc(pMsg)
	if log == "" {
		return
	}
	p.fOutErro.Println(log)
}
