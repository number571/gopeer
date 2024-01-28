package config

import (
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FLogging:  []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FLanguage: hlf_settings.CDefaultLanguage,
			FAddress: &SAddress{
				FInterface: hlf_settings.CDefaultInterfaceAddress,
				FIncoming:  hlf_settings.CDefaultIncomingAddress,
			},
			FConnection: hls_settings.CDefaultHTTPAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
