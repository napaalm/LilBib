/*
 * config.go
 *
 * File per il caricamento e gestione della configurazione
 *
 * Copyright (c) 2020 Filippo Casarin <casarin.filippo17@gmail.com>
 *
 * This file is part of LilBib.
 *
 * LilBib is free software; you can redistribute it and/or modify it
 * under the terms of the Affero GNU General Public License as
 * published by the Free Software Foundation; either version 3, or (at
 * your option) any later version.
 *
 * LilBib is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the Affero GNU General
 * Public License for more details.
 *
 * You should have received a copy of the Affero GNU General Public
 * License along with LilBib; see the file LICENSE. If not see
 * <http://www.gnu.org/licenses/>.
 */

// File per il caricamento e gestione della configurazione
package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	Generale       generale       `toml:"Generale"`
	Autenticazione autenticazione `toml:"Autenticazione"`
	LDAP           ldap           `toml:"LDAP"`
	SQL            sql            `toml:"SQL"`
}

type generale struct {
	FQDN            string `toml:"fqdn_sito"`
	Porta           string `toml:"porta_http"`
	AdminUser       string `toml:"utente_admin"`
	LunghezzaPagina uint16 `toml:"lunghezza_pagina"`
}

type autenticazione struct {
	JWTSecret     string `toml:"chiave_firma"`
	SSO           bool   `toml:"usa_sso"`
	SSOURL        string `toml:"sso_url"`
	SecureCookies bool   `toml:"cookie_sicuri"`
	DummyAuth     bool   `toml:"dummy_auth"`
}

type ldap struct {
	URI   string `toml:"uri"`
	Porta string `toml:"porta"`
}

type sql struct {
	Username  string `toml:"username"`
	Password  string `toml:"password"`
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
