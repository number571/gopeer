package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/config"
	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	cHandleIndexTemplate          = "%s" + hl_t_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "%s" + hl_t_settings.CHandleConfigSettingsPath
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHost   string
	fClient *http.Client
}

func NewRequester(pHost string, pClient *http.Client) IRequester {
	return &sRequester{
		fHost:   pHost,
		fClient: pClient,
	}
}

func (p *sRequester) GetIndex() (string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("get index (requester): %w", err)
	}

	result := string(res)
	if result != hl_t_settings.CTitlePattern {
		return "", errors.New("incorrect title pattern")
	}

	return result, nil
}

func (p *sRequester) GetSettings() (config.IConfigSettings, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get settings (requester): %w", err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.DeserializeJSON([]byte(res), cfgSettings); err != nil {
		return nil, fmt.Errorf("decode settings (requester): %w", err)
	}

	return cfgSettings, nil
}
