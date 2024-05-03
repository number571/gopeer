package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cErrorNone = iota
	cErrorGetFriends
	cErrorDecodeData
	cErrorLoadRequest
	cErrorLoadRequestID
)

func HandleNetworkRequestAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pNode anonymity.INode,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		var vRequest pkg_settings.SRequest

		if pR.Method != http.MethodPost && pR.Method != http.MethodPut {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vRequest); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		pubKey, req, errCode := unwrapRequest(pConfig, vRequest)
		switch errCode {
		case cErrorNone:
			// pass
		case cErrorGetFriends:
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: load public key")
			return
		case cErrorDecodeData:
			pLogger.PushWarn(logBuilder.WithMessage("decode_data"))
			_ = api.Response(pW, http.StatusTeapot, "failed: decode hex format data")
			return
		case cErrorLoadRequest:
			pLogger.PushWarn(logBuilder.WithMessage("load_request"))
			_ = api.Response(pW, http.StatusForbidden, "failed: decode request")
			return
		case cErrorLoadRequestID:
			pLogger.PushWarn(logBuilder.WithMessage("load_request_id"))
			_ = api.Response(pW, http.StatusForbidden, "failed: decode request id")
			return
		default:
			panic("undefined error code")
		}

		switch pR.Method {
		case http.MethodPut:
			err := pNode.SendPayload(
				pCtx,
				pubKey,
				payload.NewPayload64(uint64(pkg_settings.CServiceMask), req.ToBytes()),
			)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("send_payload"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: send payload")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: send")
			return

		case http.MethodPost:
			respBytes, err := pNode.FetchPayload(
				pCtx,
				pubKey,
				payload.NewPayload32(pkg_settings.CServiceMask, req.ToBytes()),
			)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("fetch_payload"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: fetch payload")
				return
			}

			resp, err := response.LoadResponse(respBytes)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("load_response"))
				_ = api.Response(pW, http.StatusNotExtended, "failed: load response")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, resp.ToString())
			return
		}
	}
}

func unwrapRequest(
	pConfig config.IConfig,
	pRequest pkg_settings.SRequest,
) (asymmetric.IPubKey, request.IRequest, int) {
	friends := pConfig.GetFriends()

	pubKey, ok := friends[pRequest.FReceiver]
	if !ok {
		return nil, nil, cErrorGetFriends
	}

	req := pRequest.FReqData
	if req == nil {
		return nil, nil, cErrorDecodeData
	}

	return pubKey, req, cErrorNone
}
