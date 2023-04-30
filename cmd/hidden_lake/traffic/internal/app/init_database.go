package app

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func initDatabase(pPathTo string) database.IKeyValueDB {
	return database.NewKeyValueDB(
		database.NewSettings(&database.SSettings{
			FPath:          fmt.Sprintf("%s/%s", pPathTo, hlt_settings.CPathDB),
			FLimitMessages: hlt_settings.CLimitMessages,
			FMessageSize:   hls_settings.CMessageSize,
			FWorkSize:      hls_settings.CWorkSize,
		}),
	)
}
