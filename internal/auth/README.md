# auth
Questo package implementa le funzioni per la comunicazione con il server LDAP di autenticazione e per la gestione dei token JWT di accesso.

## Funzionamento
Il server LDAP permette di autenticare gli utenti passando nome utente e password.
Una volta verificate tali credenziali viene generato un token JWT.
Esso viene passato al client, il quale poi lo include in ogni richiesta HTTP successiva.

## Server LDAP
Per comunicare con il server LDAP è necessario servirsi di questa [libreria](https://github.com/go-ldap/ldap).
Da `internal/config` va letto l'indirizzo del server.

Nota: esiste nella configurazione un'opzione `dummy_auth` la quale fa sì che il server accetti ogni utente senza controllare le credenziali.

## Token JWT
I token JWT conservano il nome utente e il livello di privilegi dell'utente (user o admin).
Essi sono firmati crittograficamente, e l'algoritmo scelto per questo progetto è HMAC-SHA256.
Il secret per la firma dev'essere ottenuto da `internal/config`.
[Libreria](https://github.com/gbrlsnchs/jwt).

## Tipi
```go
type UserInfo
```

## Errori
```go
type AuthenticationError
type LDAPError
type JWTCreationError
type InvalidTokenError
```

## Funzioni
```go
func InitializeSigning()
func AuthenticateUser(username, password string) ([]byte, error) // ritorna il token
func ParseToken(token []byte) (UserInfo, error)
```
