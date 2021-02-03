/*
 * auth_test.go
 *
 * File di test per il package auth.
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

import (
	"fmt"
	"git.napaalm.xyz/napaalm/LilBib/internal/config"
	"testing"
)

func TestAll(t *testing.T) {
	// Carica e visualizza la configurazione di test
	config.LoadConfig("./config_test.toml")
	fmt.Println(config.Config)

	// Genera un token
	InitializeSigning()
	token, err := AuthenticateUser("admin", "password")

	if err == nil {
		fmt.Println(string(token))
	} else {
		fmt.Println(err)
	}

	// Verifica il token generato
	userInfo, err := ParseToken(token)

	if err == nil {
		fmt.Println(userInfo)
	} else {
		fmt.Println(err)
	}

	userInfo, err = ParseToken([]byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXlsb2FkIjp7ImlzcyI6IkxpbEJpYiIsInN1YiI6ImFkbWluIiwiYXVkIjpbImh0dHA6Ly9leGFtcGxlLm9yZyIsImh0dHBzOi8vZXhhbXBsZS5vcmciXSwiZXhwIjoxNTkzNTQwNDA5LCJpYXQiOjE1OTA5NDg0MDl9LCJpc0FkbWluIjpmYWxzZX0.0bjvYV3d37Rlyxe9Mof1hMp37nb31V6NT1EpOijwBVA"))

	if err == nil {
		fmt.Println(userInfo)
	} else {
		fmt.Println(err)
	}
}
