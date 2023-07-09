package config

import "github.com/number571/go-peer/pkg/logger"

type IConfig interface {
	GetLogging() logger.ILogging
	GetNetwork() string
	GetAddress() IAddress
	GetConnections() []string
	GetConsumers() []string
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
}
