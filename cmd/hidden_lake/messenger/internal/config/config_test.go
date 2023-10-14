package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/pkg/filesystem"
)

const (
	tcLogging    = true
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `{
	"settings": {
		"messages_capacity": %d,
		"work_size_bits": %d
	},
	"logging": ["info", "erro"],
	"language": "RUS",
	"address": {
		"interface": "%s",
		"incoming": "%s",
		"pprof": "%s"
	},
	"connection": "%s",
	"backup_connections": [
		"%s"
	],
	"storage_key": "%s"
}`
)

const (
	tcAddressInterface  = "address_interface"
	tcAddressIncoming   = "address_incoming"
	tcAddressPPROF      = "address_pprof"
	tcConnectionService = "connection_service"
	tcConnectionBackup  = "connection_backup"
	tcStorageKey        = "storage_key"
	tcMessageSize       = (1 << 20)
	tcWorkSize          = 20
	tcKeySize           = 1024
	tcMessagesCapacity  = 1000
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessagesCapacity,
		tcWorkSize,
		tcAddressInterface,
		tcAddressIncoming,
		tcAddressPPROF,
		tcConnectionService,
		tcConnectionBackup,
		tcStorageKey,
	)
}

func testConfigDefaultInit(configPath string) {
	filesystem.OpenFile(configPath).Write([]byte(testNewConfigString()))
}

func TestConfig(t *testing.T) {
	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.GetSettings().GetMessagesCapacity() != tcMessagesCapacity {
		t.Error("settings key size is invalid")
		return
	}

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSize {
		t.Error("settings key size is invalid")
		return
	}

	if cfg.GetLogging().HasInfo() != tcLogging {
		t.Error("logging.info is invalid")
		return
	}

	if cfg.GetLogging().HasErro() != tcLogging {
		t.Error("logging.erro is invalid")
		return
	}

	if cfg.GetLogging().HasWarn() == tcLogging {
		t.Error("logging.warn is invalid")
		return
	}

	if cfg.GetLanguage() != utils.CLangRUS {
		t.Error("language is invalid")
		return
	}

	if cfg.GetAddress().GetInterface() != tcAddressInterface {
		t.Error("address.interface is invalid")
		return
	}

	if cfg.GetAddress().GetIncoming() != tcAddressIncoming {
		t.Error("address.incoming is invalid")
		return
	}

	if cfg.GetAddress().GetPPROF() != tcAddressPPROF {
		t.Error("address.pprof is invalid")
		return
	}

	if cfg.GetConnection() != tcConnectionService {
		t.Error("connection.service is invalid")
		return
	}

	if len(cfg.GetBackupConnections()) != 1 {
		t.Error("length of connections.backup is invalid")
		return
	}

	if cfg.GetBackupConnections()[0] != tcConnectionBackup {
		t.Error("connections[0].backup is invalid")
		return
	}

	if cfg.GetStorageKey() != tcStorageKey {
		t.Error("storage_key is invalid")
		return
	}
}
