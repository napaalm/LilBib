/*
 * auth.go
 *
 * Funzione per autenticare un utente.
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

package auth

import "git.antonionapolitano.eu/napaalm/LilBib/internal/config"

// Informazioni sull'utente
type UserInfo struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Group    string `json:"group"`
}

// Verifica le credenziali e restituisce il token.
func AuthenticateUser(username, password string) ([]byte, error) {

	var (
		err      error
		userInfo UserInfo
	)

	// Controlla le credenziali
	if !config.Config.Autenticazione.DummyAuth {
		userInfo, err = checkCredentials(username, password)
	} else {
		err = nil
		userInfo = UserInfo{username, "1337 h4x0r", "h4x0rz"}
	}

	if err == nil {
		// Genera il token
		token, err := getToken(userInfo)

		if err != nil {
			return nil, err
		}

		return token, nil

	} else {
		return nil, err
	}
}
