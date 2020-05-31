# hash
Pacchetto per la generazione e verifica dei codici dei libri.

## Funzionamento

password: 4 byte per il codice del libro e 16 per la password effettiva  
hash: sha256(password effettiva)

## Funzioni
```go
func Verifica(password []byte) (db.Libro, error)
func Genera(codice uint32) ([]byte, []byte, error) // hash password
```

