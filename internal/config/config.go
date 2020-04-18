package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	Generale generale `toml:"Generale"`
	LDAP     ldap     `toml:"LDAP"`
	SQL      sql      `toml:"SQL"`
}

type generale struct {
	AdminUser string `toml:"admin"`
	JWTSecret string `toml:"chiave_firma"`
}

type ldap struct {
	URI   string `toml:"uri"`
	Porta string `toml:"porta"`
}

type sql struct {
	Username  string `toml:"username"`
	Password  string `toml:"passowrd"`
	Indirizzo string `toml:"indirizzo"`
	Database  string `toml:"database"`
}

var Config config

func LoadConfig(path string) error {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	if _, err := toml.DecodeFile(absPath, &Config); err != nil {
		return err
	}

	return nil
}
