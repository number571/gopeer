package handler

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	hlm_client "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/web"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

type sChatMessage struct {
	FIsIncoming bool
	utils.SMessage
}

type sChatAddress struct {
	FAliasName  string
	FPubKeyHash string
}

type sChatMessages struct {
	*sTemplate
	FAddress  sChatAddress
	FMessages []sChatMessage
}

func FriendsChatPage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pDB database.IKVDatabase,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/friends/chat" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		_ = pR.ParseForm()
		_ = pR.ParseMultipartForm(10 << 20) // default value

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			ErrorPage(pLogger, pCfg, "get_alias_name", "alias name is nil")(pW, pR)
			return
		}

		client := getHLSClient(pCfg)
		myPubKey, err := client.GetPubKey(pCtx)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_public_key", "read public key")(pW, pR)
			return
		}

		recvPubKey, err := getReceiverPubKey(pCtx, client, aliasName)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_receiver", "get receiver by public key")(pW, pR)
			return
		}

		rel := database.NewRelation(myPubKey, recvPubKey)

		switch pR.FormValue("method") {
		case http.MethodPost, http.MethodPut:
			msgBytes, err := getMessageBytes(pR)
			if err != nil {
				ErrorPage(pLogger, pCfg, "get_message", "get message bytes")(pW, pR)
				return
			}

			if err := sendMessage(pCtx, pCfg, client, aliasName, msgBytes); err != nil {
				ErrorPage(pLogger, pCfg, "send_message", "push message to network")(pW, pR)
				return
			}

			dbMsg := database.NewMessage(false, pCfg.GetSettings().GetPseudonym(), msgBytes)
			if err := pDB.Push(rel, dbMsg); err != nil {
				ErrorPage(pLogger, pCfg, "push_message", "add message to database")(pW, pR)
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
			http.Redirect(pW, pR,
				fmt.Sprintf("/friends/chat?alias_name=%s", aliasName),
				http.StatusSeeOther)
			return
		}

		start := uint64(0)
		size := pDB.Size(rel)

		messagesCap := pCfg.GetSettings().GetMessagesCapacity()
		if size > messagesCap {
			start = size - messagesCap
		}

		msgs, err := pDB.Load(rel, start, size)
		if err != nil {
			ErrorPage(pLogger, pCfg, "read_database", "read database")(pW, pR)
			return
		}

		res := &sChatMessages{
			sTemplate: getTemplate(pCfg),
			FAddress: sChatAddress{
				FAliasName:  aliasName,
				FPubKeyHash: recvPubKey.GetHasher().ToString(),
			},
			FMessages: make([]sChatMessage, 0, len(msgs)),
		}

		for _, msg := range msgs {
			res.FMessages = append(res.FMessages, sChatMessage{
				FIsIncoming: msg.IsIncoming(),
				SMessage: getMessage(
					false,
					msg.GetPseudonym(),
					msg.GetMessage(),
					msg.GetTimestamp(),
				),
			})
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"chat.html",
		)
		if err != nil {
			panic(err)
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = t.Execute(pW, res)
	}
}

func getMessageBytes(pR *http.Request) ([]byte, error) {
	switch pR.FormValue("method") {
	case http.MethodPost:
		strMsg := strings.TrimSpace(pR.FormValue("input_message"))
		if strMsg == "" {
			return nil, errors.New("error: message is nil")
		}
		if utils.HasNotWritableCharacters(strMsg) {
			return nil, errors.New("error: message has not writable characters")
		}
		return wrapText(strMsg), nil
	case http.MethodPut:
		filename, fileBytes, err := getUploadFile(pR)
		if err != nil {
			return nil, fmt.Errorf("error: upload file: %w", err)
		}
		return wrapFile(filename, fileBytes), nil
	default:
		panic("got not supported method")
	}
}

func getUploadFile(pR *http.Request) (string, []byte, error) {
	// Get handler for filename, size and headers
	file, handler, err := pR.FormFile("input_file")
	if err != nil {
		return "", nil, fmt.Errorf("error: receive file: %w", err)
	}
	defer file.Close()

	if handler.Size == 0 {
		return "", nil, errors.New("error: file size is nil")
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", nil, fmt.Errorf("error: read file bytes: %w", err)
	}

	return handler.Filename, fileBytes, nil
}

func sendMessage(
	pCtx context.Context,
	pCfg config.IConfig,
	pClient hls_client.IClient,
	pAliasName string,
	pMsgBytes []byte,
) error {
	msgLimit, err := utils.GetMessageLimit(pCtx, pClient)
	if err != nil {
		return fmt.Errorf("error: try send message: %w", err)
	}

	if uint64(len(pMsgBytes)) > (msgLimit + symmetric.CAESBlockSize + hashing.CSHA256Size) {
		return fmt.Errorf("error: len message > limit: %w", err)
	}

	hlmClient := hlm_client.NewClient(
		hlm_client.NewBuilder(),
		hlm_client.NewRequester(pClient),
	)
	return hlmClient.PushMessage(pCtx, pAliasName, pCfg.GetSettings().GetPseudonym(), pMsgBytes)
}

func getReceiverPubKey(
	pCtx context.Context,
	client hls_client.IClient,
	aliasName string,
) (asymmetric.IPubKey, error) {
	friends, err := client.GetFriends(pCtx)
	if err != nil {
		return nil, fmt.Errorf("error: read friends: %w", err)
	}

	friendPubKey, ok := friends[aliasName]
	if !ok {
		return nil, errors.New("undefined public key by alias name")
	}

	return friendPubKey, nil
}
