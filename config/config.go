package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type Server struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`

	NickservPassword string `json:"nickserv_password"`

	Channels []string `json:"channels"`
}

type Config struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Realname string `json:"realname"`

	Admins []string `json:"admins"`

	Server Server `json:"server"`

	Debug bool `json:"debug"`
}

func (c *Config) Check() error {
	if len(c.Nickname) < 3 {
		return errors.New("Nickname is empty or too short")
	}

	if c.Username == "" {
		return errors.New("Username can't be empty")
	}

	if c.Realname == "" {
		return errors.New("Realname can't be empty")
	}

	if c.Server.Address == "" {
		return errors.New("Server address not specified")
	}

	if c.Server.Port == 0 {
		return errors.New("Server port can't be zero")
	}

	return nil
}

var conf *Config

func Parse(s string) (*Config, error) {
	f, err := os.Open(s)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening the config file")
	}
	defer f.Close()

	conf = new(Config)
	err = json.NewDecoder(f).Decode(conf)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding the config file")
	}
	return conf, nil
}

func Get() *Config {
	if conf == nil {
		return new(Config)
	}
	return conf
}
