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

// Verifica le credenziali, ottiene il livello di permessi dell'utente e restituisce il token.
func AuthenticateUser(username, password string) ([]byte, error) {

	// Controlla le credenziali
	var err error
	if !config.Config.Autenticazione.DummyAuth {
		err = checkCredentials(username, password)
	} else {
		err = nil
	}

	if err == nil {
		// Controlla se l'utente è admin
		isAdmin := (username == config.Config.Generale.AdminUser)

		// Genera il token
		token, err := getToken(username, isAdmin)

		if err != nil {
			return nil, err
		}

		return token, nil

	} else {
		return nil, err
	}
}
