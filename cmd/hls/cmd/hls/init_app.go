package main

import (
	"flag"
	"fmt"

	"github.com/number571/go-peer/cmd/hls/internal/app"
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/filesystem"
)

func initApp() (modules.IApp, error) {
	var (
		inputKey string
	)

	flag.StringVar(&inputKey, "key", "priv.key", "input private key from file")
	flag.Parse()

	privKeyStr, err := filesystem.OpenFile(inputKey).Read()
	if err != nil {
		return nil, err
	}

	privKey := asymmetric.LoadRSAPrivKey(string(privKeyStr))
	if privKey == nil {
		return nil, fmt.Errorf("private key is invalid")
	}

	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}

	return app.NewApp(cfg, initNode(cfg, privKey)), nil
}
