package config

import (
	"fmt"
	"os"
	"testing"

	hla_settings "github.com/number571/go-peer/cmd/hidden_lake/adapters/common/producer/pkg/settings"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

func TestInit(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 1)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)

	config1, err := InitConfig(configFile, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if config1.GetConnection().GetSrvHost() != tcAddress1 {
		t.Error("got invalid field with exist config (1)")
		return
	}

	os.Remove(configFile)
	if err := os.WriteFile(configFile, []byte("abc"), 0o600); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitConfig(configFile, nil); err == nil {
		t.Error("success init config with invalid config structure (1)")
		return
	}

	os.Remove(configFile)

	if _, err := InitConfig(configFile, &SConfig{}); err == nil {
		t.Error("success init config with invalid config structure (2)")
		return
	}

	os.Remove(configFile)

	config2, err := InitConfig(configFile, config1.(*SConfig))
	if err != nil {
		t.Error(err)
		return
	}

	if config2.GetConnection().GetSrvHost() != tcAddress1 {
		t.Error("got invalid field with exist config (2)")
		return
	}

	os.Remove(configFile)

	config3, err := InitConfig(configFile, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if config3.GetConnection().GetSrvHost() != hla_settings.CDefaultSrvAddress {
		t.Error("got invalid field with exist config (3)")
		return
	}

	if config3.GetAddress() != hla_settings.CDefaultHTTPAddress {
		t.Error("got invalid field with exist config (3)")
		return
	}
}
