package config

import "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
