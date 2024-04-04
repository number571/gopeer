package handler

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/internal/language"
)

type sTemplate struct {
	FLanguage language.ILanguage
}

func getTemplate(pCfg config.IConfig) *sTemplate {
	return &sTemplate{
		FLanguage: pCfg.GetSettings().GetLanguage(),
	}
}

func getHLSClient(pCfg config.IConfig) hls_client.IClient {
	return hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+pCfg.GetConnection(),
			&http.Client{Timeout: (10 * time.Minute)},
		),
	)
}
