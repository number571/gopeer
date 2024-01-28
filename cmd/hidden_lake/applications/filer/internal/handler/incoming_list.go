package handler

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/internal/config"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigListHTTP(pLogger logger.ILogger, pCfg config.IConfig, pPathTo string) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		fPubKey := asymmetric.LoadRSAPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			panic("public key is nil (invalid data from HLS)!")
		}

		page, err := strconv.Atoi(pR.URL.Query().Get("page"))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_page"))
			api.Response(pW, http.StatusBadRequest, "failed: incorrect page")
			return
		}

		result, err := getListFileInfo(pCfg, pPathTo, page)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("open storage"))
			api.Response(pW, http.StatusInternalServerError, "failed: open storage")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, result)
	}
}

func getListFileInfo(pCfg config.IConfig, pPathTo string, pPage int) ([]hlf_settings.SFileInfo, error) {
	entries, err := os.ReadDir(hlf_settings.CPathSTG)
	if err != nil {
		return nil, err
	}
	result := make([]hlf_settings.SFileInfo, 0, len(entries))
	for i := (pPage * hlf_settings.CPageOffset); i < len(entries); i++ {
		e := entries[i]
		if e.IsDir() {
			continue
		}
		if i != (pPage*hlf_settings.CPageOffset) && i%hlf_settings.CPageOffset == 0 {
			break
		}
		fullPath := fmt.Sprintf("%s/%s/%s", pPathTo, hlf_settings.CPathSTG, e.Name())
		result = append(result, hlf_settings.SFileInfo{
			FName: e.Name(),
			FHash: getFileHash(fullPath),
			FSize: getFileSize(fullPath),
		})
	}
	return result, nil
}

func getFileSize(filename string) uint64 {
	stat, _ := os.Stat(filename)
	return uint64(stat.Size())
}

func getFileHash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}

	return encoding.HexEncode(h.Sum(nil))
}
