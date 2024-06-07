package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	ConfigPath       string                    `toml:"config_path"`
	InstallationPath string                    `toml:"installation_path"`
	LogLevel         string                    `toml:"log_level"`
	MaxNoops         int                       `toml:"max_noops"`
	Scrapers         map[string]ScraperConfigs `toml:"scrapers"`
}
type ScraperConfigs = []ScraperConfig
type ScraperConfig = map[string]any

var ProdexConf Config

func LoadConfig(conf string) error {
	var err error
	if conf == "" {
		_, err = toml.DecodeFile("dev/prodex.toml", &ProdexConf)
	} else {
		_, err = toml.Decode(conf, &ProdexConf)
	}
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	return nil
}
