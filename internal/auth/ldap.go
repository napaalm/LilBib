/*
 * ldap.go
 *
 * Autenticazione su un server LDAP.
 *
 * Copyright (c) 2020 Antonio Napolitano <nap@antonionapolitano.eu>
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

// Package per l'autenticazione degli utenti.
package auth

import (
	"fmt"

	"git.antonionapolitano.eu/napaalm/LilBib/internal/config"
	. "github.com/go-ldap/ldap/v3"
)

// Errore di autenticazione
type AuthenticationError struct {
	user string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error for user %s.", e)
}

// Controlla le credenziali sul server LDAP
func checkCredentials(username string, password string) error {

	// Utente readonly per la ricerca dell'utente effettivo
	bindusername := "readonly"
	bindpassword := "password"

	// Ottiene l'indirizzo del server dalla configurazione
	host := config.Config.LDAP.URI
	port := config.Config.LDAP.Porta

	// Connessione al server LDAP
	l, err := DialURL("ldap://" + host + ":" + port)
	if err != nil {
		return &AuthenticationError{username}
	}
	defer l.Close()

	// Per prima cosa effettuo l'accesso con un utente readonly
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		return &AuthenticationError{username}
	}

	// Cerco l'username richiesto
	searchRequest := NewSearchRequest(
		"dc=example,dc=com",
		ScopeWholeSubtree, NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return &AuthenticationError{username}
	}

	// Verifico il numero di utenti corrispondenti e ottendo il DN dell'utente
	if len(sr.Entries) != 1 {
		return &AuthenticationError{username}
	}

	userdn := sr.Entries[0].DN

	// Verifica la password
	err = l.Bind(userdn, password)
	if err != nil {
		return &AuthenticationError{username}
	}

	return nil
}
