package config

import (
	_ "embed"

	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util/fs"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var DefaultConfig []byte

const configFile = "config.yaml"

type Config struct {
	Account struct {
		Uin      int64  `yaml:"uid"`
		Password string `yaml:"password"`
	} `yaml:"account"`
}

func Read() *Config {
	if !fs.FileExists(configFile) {
		log.Error().Msg("配置文件不存在")
		log.Info().Msgf("已生成配置文件模板 %s", configFile)
		fs.MustWriteFile(configFile, DefaultConfig, 0644)
		return nil
	}
	fd, err := fs.OpenFile(configFile)
	if err != nil {
		log.Error().Err(err).Msg("无法读取配置文件")
		return nil
	}
	c := new(Config)
	if err = yaml.NewDecoder(fd).Decode(c); err != nil {
		log.Error().Err(err).Msg("配置文件格式不正确")
		return nil
	}

	if c.Account.Uin == 0 || len(c.Account.Password) == 0 {
		log.Error().Msg("账号密码未设置")
		return nil
	}

	return c
}
