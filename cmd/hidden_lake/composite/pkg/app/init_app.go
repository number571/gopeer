package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/internal/flag"
	"github.com/number571/go-peer/pkg/types"

	hlf_app "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/app"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	hlm_app "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/app"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/composite/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/composite/pkg/settings"
	hle_app "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/app"
	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	hll_app "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/app"
	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/settings"
	hlt_app "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/app"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
	hls_app "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/app"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	servicesCount = 6
)

func InitApp(pArgs []string, pDefaultPath, pDefaultKey string, pParallel uint64) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfg, err := config.InitConfig(filepath.Join(inputPath, settings.CPathYML), nil)
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	runners, err := getRunners(cfg, pArgs, pDefaultPath, pDefaultKey, pParallel)
	if err != nil {
		return nil, fmt.Errorf("get runners: %w", err)
	}

	return NewApp(cfg, runners), nil
}

func getRunners(pCfg config.IConfig, pArgs []string, pDefaultPath, pDefaultKey string, pParallel uint64) ([]types.IRunner, error) {
	runners := make([]types.IRunner, 0, servicesCount)

	var (
		runner types.IRunner
		err    error
	)

	for _, sName := range pCfg.GetServices() {
		switch sName {
		case hls_settings.CServiceFullName:
			runner, err = hls_app.InitApp(pArgs, pDefaultPath, pDefaultKey, pParallel)
		case hlt_settings.CServiceFullName:
			runner, err = hlt_app.InitApp(pArgs, pDefaultPath)
		case hle_settings.CServiceFullName:
			runner, err = hle_app.InitApp(pArgs, pDefaultPath, pDefaultKey, pParallel)
		case hll_settings.CServiceFullName:
			runner, err = hll_app.InitApp(pArgs, pDefaultPath)
		case hlm_settings.CServiceFullName:
			runner, err = hlm_app.InitApp(pArgs, pDefaultPath)
		case hlf_settings.CServiceFullName:
			runner, err = hlf_app.InitApp(pArgs, pDefaultPath)
		default:
			return nil, fmt.Errorf("unknown service %s", sName)
		}
		if err != nil {
			return nil, err
		}
		runners = append(runners, runner)
	}

	return runners, nil
}
