package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Cloudfront CloudfrontConfig `mapstructure:"cloudfront"`
	TLS        TLSConfig        `mapstructure:"tls"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type CloudfrontConfig struct {
	URL        string `mapstructure:"url"`
	KeyID      string `mapstructure:"key_id"`
	PrivateKey string `mapstructure:"private_key"`
}

type TLSConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
	CAFile     string `mapstructure:"ca_file"`
	ClientAuth string `mapstructure:"client_auth"` // "none", "request", "require", "verify-if-given", "require-and-verify"
	MinVersion string `mapstructure:"min_version"` // "1.0", "1.1", "1.2", "1.3"
	MaxVersion string `mapstructure:"max_version"` // "1.0", "1.1", "1.2", "1.3"
}

func LoadConfig(path string) (config Config, err error) {
	// parse the provided input
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Config{}, err
	}
	base := filepath.Base(absPath)
	baseDir := filepath.Dir(absPath)
	fileExt := filepath.Ext(base)
	fileName := base[:len(base)-len(fileExt)]

	// add config to viper
	viper.AddConfigPath(baseDir)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")

	// allow overridding from env
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	return
}
