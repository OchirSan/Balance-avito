package models

import (
	"io/ioutil"
	"strings"
	"time"
	"github.com/BurntSushi/toml"
)

var (
	configPath = "config/config.toml"
	hashPaths  = []string{configPath}
)

type duration time.Duration

// Config struct
type Config struct {
	SQLDataBase SQLDataBase `toml:"SQLDataBase"`
	ServerOpt   ServerOpt   `toml:"ServerOpt"`
	HashSum     []byte
}

func (d *duration) UnmarshalText(text []byte) error {
	temp, err := time.ParseDuration(string(text))
	*d = duration(temp)
	return err
}

// ServerOpt struct
type ServerOpt struct {
	ReadTimeout  duration
	WriteTimeout duration
	IdleTimeout  duration
}

// LoadConfig from path
func LoadConfig(c *Config) {
	_, err := toml.DecodeFile(configPath, c)
	if err != nil {
		return
	}
	c.SQLDataBase.UserID = getCredential("../etc/scrt/balance/sqlUser.txt")
	c.SQLDataBase.Password = getCredential("../etc/scrt/balance/sqlPassword.txt")

}

func getCredential(path string) string {
	hashPaths = append(hashPaths, path)
	c, _ := ioutil.ReadFile(path)
	return strings.TrimSpace(string(c))
}