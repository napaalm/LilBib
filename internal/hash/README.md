# hash
Pacchetto per la generazione e verifica dei codici dei libri.

## Funzionamento

password: `base64_urlsafe(4 byte per il codice del libro e 16 per la password effettiva)`
hash: `base64(sha256(password effettiva))`

## Funzioni
```go
func Verifica(pass string) (db.Libro, error)
func Genera(codice uint32) (string, error)
```

