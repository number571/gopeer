package config

import (
	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/internal/settings"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/filesystem"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if filesystem.OpenFile(cfgPath).IsExist() {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, errors.WrapError(err, "load config")
		}
		return cfg, nil
	}
	if initCfg == nil {
		initCfg = &SConfig{
			SConfigSettings: settings.SConfigSettings{
				FSettings: settings.SConfigSettingsBlock{
					FMessageSizeBytes:   hls_settings.CDefaultMessageSize,
					FWorkSizeBits:       hls_settings.CDefaultWorkSize,
					FQueuePeriodMS:      hls_settings.CDefaultQueuePeriod,
					FLimitVoidSizeBytes: hls_settings.CDefaultLimitVoidSize,
				},
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  hlt_settings.CDefaultTCPAddress,
				FHTTP: hlt_settings.CDefaultHTTPAddress,
			},
			FConnections: []string{
				hlt_settings.CDefaultConnectionAddress,
			},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, errors.WrapError(err, "build config")
	}
	return cfg, nil
}
