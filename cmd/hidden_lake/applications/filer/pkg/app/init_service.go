package app

import (
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/internal/handler"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/web"
	"github.com/number571/go-peer/pkg/logger"
)

func (p *sApp) initIncomingServiceHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlm_settings.CLoadPath,
		handler.HandleIncomigLoadHTTP(p.fHTTPLogger, p.fConfig, p.fPathTo),
	) // POST

	mux.HandleFunc(
		hlm_settings.CListPath,
		handler.HandleIncomigListHTTP(p.fHTTPLogger, p.fConfig, p.fPathTo),
	) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetIncoming(),
		ReadTimeout: (5 * time.Second),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}
}

func (p *sApp) initInterfaceServiceHTTP() {
	mux := http.NewServeMux()
	mux.Handle(hlm_settings.CStaticPath, http.StripPrefix(
		hlm_settings.CStaticPath,
		handleFileServer(p.fHTTPLogger, p.fConfig, http.FS(web.GetStaticPath()))),
	)

	cfgWrapper := config.NewWrapper(p.fConfig)

	mux.HandleFunc(hlm_settings.CHandleIndexPath, handler.IndexPage(p.fHTTPLogger, p.fConfig))                       // GET, POST
	mux.HandleFunc(hlm_settings.CHandleFaviconPath, handler.FaviconPage(p.fHTTPLogger, p.fConfig))                   // GET
	mux.HandleFunc(hlm_settings.CHandleAboutPath, handler.AboutPage(p.fHTTPLogger, p.fConfig))                       // GET
	mux.HandleFunc(hlm_settings.CHandleSettingsPath, handler.SettingsPage(p.fHTTPLogger, cfgWrapper))                // GET, PATCH, PUT, POST, DELETE
	mux.HandleFunc(hlm_settings.CHandleFriendsPath, handler.FriendsPage(p.fHTTPLogger, p.fConfig))                   // GET, POST, DELETE
	mux.HandleFunc(hlm_settings.CHandleFriendsStoragePath, handler.StoragePage(p.fHTTPLogger, p.fConfig, p.fPathTo)) // GET, POST, DELETE

	p.fIntServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetInterface(),
		ReadTimeout: (5 * time.Second),
		Handler:     mux, // http.TimeoutHandler send panic from websocket use
	}
}

func handleFileServer(pLogger logger.ILogger, pCfg config.IConfig, pFS http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := pFS.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(pLogger, pCfg)(w, r)
			return
		}
		http.FileServer(pFS).ServeHTTP(w, r)
	})
}
