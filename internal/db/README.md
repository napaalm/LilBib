# db
Pacchetto per comunicare con il database SQL

## Tipi
```go
type Libro
type Autore
type Genere
type Prestito
```

## Funzioni
```go
func GetLibro(codice uint32) (Libro, error)
func GetLibri(page uint32) ([]Libri, error)
func GetAutori(iniziale uint8) ([]Autore, error)
func GetGeneri() ([]Genere, error)
func GetPrestiti(utente string) ([]Prestito, error)

func RicercaLibri(nome string, autore, genere []uint32) ([]Libro, error)
func RicercaAutori(nome string) ([]Autore, error)
func RicercaGeneri(nome string) ([]Genere, error)

func AddLibro(titolo string, autore, genere uint32, hash string) error
func AddGenere(nome string) error
func AddAutore(nome, cognome string) error
func AddPrestito(codice, libro uint32, utente string data_prenotazione TipoData?, durata uint32) error

func SetRestituito(prestito uint32, data_restituzione TipoData?) error
```

