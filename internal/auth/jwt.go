/*
 * jwt.go
 *
 * Funzioni per la generazione dei token JWT di autenticazione e per
 * la loro verifica.
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
	"time"

	"git.napaalm.xyz/napaalm/LilBib/internal/config"
	"github.com/gbrlsnchs/jwt/v3"
)

var (
	jwtSigner *jwt.HMACSHA
)

// Errore di generazione del token
type JWTCreationError struct {
	username string
}

func (e *JWTCreationError) Error() string {
	return fmt.Sprintf("Failed to sign the JWT token for username %s.", e.username)
}

// Errore nella verifica del token
type InvalidTokenError struct {
	token []byte
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("Failed to verify the following token: %s", string(e.token))
}

// Formato del payload JWT
type customPayload struct {
	Payload  jwt.Payload
	FullName string `json:"full_name"`
	Group    string `json:"group"`
}

var dummyUserInfo = UserInfo{
	"h4x0r",
	"1337 h4x0r",
	"1337",
}

// Inizializza l'algoritmo per la firma HS256
func InitializeSigning() {

	// Ottiene la chiave segreta dalla configurazione
	secret := config.Config.Autenticazione.JWTSecret

	// Inizializza l'algoritmo
	jwtSigner = jwt.NewHS256([]byte(secret))
}

// Genera un token
func getToken(userInfo UserInfo) ([]byte, error) {

	// Ottiene il tempo corrente
	now := time.Now()

	// Copio i valori necessari dalla configurazione
	fqdn := config.Config.Generale.FQDN

	// Definisco il payload
	pl := customPayload{
		Payload: jwt.Payload{
			Issuer:         "LilBib",
			Subject:        userInfo.Username,
			Audience:       jwt.Audience{"http://" + fqdn, "https://" + fqdn},
			ExpirationTime: jwt.NumericDate(now.Add(24 * time.Hour)),
			IssuedAt:       jwt.NumericDate(now),
		},
		FullName: userInfo.FullName,
		Group:    userInfo.Group,
	}

	// Firma il token
	token, err := jwt.Sign(pl, jwtSigner)

	if err != nil {
		return nil, &JWTCreationError{userInfo.Username}
	}

	return token, nil
}

// Verifica un token e ne restituisce le informazioni
func ParseToken(token []byte) (UserInfo, error) {

	var (
		// Ottengo il tempo corrente
		now = time.Now()

		// Carico l'FQDN dalla configurazione e definisco l'audience
		fqdn = config.Config.Generale.FQDN
		aud  = jwt.Audience{"http://" + fqdn, "https://" + fqdn}

		// Inizializzo i "validatori"
		iatValidator = jwt.IssuedAtValidator(now)
		expValidator = jwt.ExpirationTimeValidator(now)
		audValidator = jwt.AudienceValidator(aud)

		// Costruisco il validatore supremo
		pl              customPayload
		validatePayload = jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator, audValidator)
	)

	// Verifico il token
	_, err := jwt.Verify(token, jwtSigner, &pl, validatePayload)

	if err != nil {
		// Valori di errore
		return dummyUserInfo, &InvalidTokenError{token}
	}

	// Ottengo le informazioni sull'utente
	username := pl.Payload.Subject
	fullName := pl.FullName
	group := pl.Group

	return UserInfo{username, fullName, group}, nil
}
