package config

import (
	"errors"
	"fmt"
	"os"

	logger "github.com/number571/go-peer/internal/logger/std"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
	_ logger.ILogging = &sLogging{}
)

type SConfigSettings struct {
	FMessageSizeBytes   uint64 `yaml:"message_size_bytes"`
	FWorkSizeBits       uint64 `yaml:"work_size_bits"`
	FKeySizeBits        uint64 `yaml:"key_size_bits"`
	FMessagesCapacity   uint64 `yaml:"messages_capacity"`
	FQueuePeriodMS      uint64 `yaml:"queue_period_ms,omitempty"`
	FLimitVoidSizeBytes uint64 `yaml:"limit_void_size_bytes,omitempty"`
	FNetworkKey         string `yaml:"network_key,omitempty"`
}

type SConfig struct {
	FSettings *SConfigSettings `yaml:"settings"`

	FLogging     []string  `yaml:"logging,omitempty"`
	FAddress     *SAddress `yaml:"address,omitempty"`
	FStorage     bool      `yaml:"storage,omitempty"`
	FConnections []string  `yaml:"connections,omitempty"`
	FConsumers   []string  `yaml:"consumers,omitempty"`

	fLogging *sLogging
}

type SAddress struct {
	FTCP   string `yaml:"tcp,omitempty"`
	FHTTP  string `yaml:"http,omitempty"`
	FPPROF string `yaml:"pprof,omitempty"`
}

type sLogging []bool

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pFilepath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("config file '%s' already exist", pFilepath)
	}

	if err := pCfg.initConfig(); err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	if err := os.WriteFile(pFilepath, encoding.SerializeYAML(pCfg), 0o644); err != nil {
		return nil, fmt.Errorf("write config: %w", err)
	}

	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	if _, err := os.Stat(pFilepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file '%s' does not exist", pFilepath)
	}

	bytes, err := os.ReadFile(pFilepath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := new(SConfig)
	if err := encoding.DeserializeYAML(bytes, cfg); err != nil {
		return nil, fmt.Errorf("deserialize config: %w", err)
	}

	if err := cfg.initConfig(); err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}
	return cfg, nil
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) GetStorage() bool {
	return p.FStorage
}

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetQueuePeriodMS() uint64 {
	return p.FQueuePeriodMS
}

func (p *SConfigSettings) GetMessagesCapacity() uint64 {
	return p.FMessagesCapacity
}

func (p *SConfigSettings) GetLimitVoidSizeBytes() uint64 {
	return p.FLimitVoidSizeBytes
}

func (p *SConfigSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfigSettings) GetKeySizeBits() uint64 {
	return p.FKeySizeBits
}

func (p *SConfig) isValid() bool {
	return true &&
		p.FSettings.FMessageSizeBytes != 0 &&
		p.FSettings.FMessagesCapacity != 0 &&
		p.FSettings.FKeySizeBits != 0
}

func (p *SConfig) initConfig() error {
	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
	}

	if p.FAddress == nil {
		p.FAddress = new(SAddress)
	}

	if !p.isValid() {
		return errors.New("load config settings")
	}

	if err := p.loadLogging(); err != nil {
		return fmt.Errorf("load logging: %w", err)
	}

	return nil
}

func (p *SConfig) loadLogging() error {
	// [info, warn, erro]
	logging := sLogging(make([]bool, 3))

	mapping := map[string]int{
		"info": 0,
		"warn": 1,
		"erro": 2,
	}

	for _, v := range p.FLogging {
		logType, ok := mapping[v]
		if !ok {
			return fmt.Errorf("undefined log type '%s'", v)
		}
		logging[logType] = true
	}

	p.fLogging = &logging
	return nil
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SAddress) GetTCP() string {
	return p.FTCP
}

func (p *SAddress) GetHTTP() string {
	return p.FHTTP
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}

func (p *SConfig) GetConnections() []string {
	return p.FConnections
}

func (p *SConfig) GetConsumers() []string {
	return p.FConsumers
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}

func (p *sLogging) HasInfo() bool {
	return (*p)[0]
}

func (p *sLogging) HasWarn() bool {
	return (*p)[1]
}

func (p *sLogging) HasErro() bool {
	return (*p)[2]
}
