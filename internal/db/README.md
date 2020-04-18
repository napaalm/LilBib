# db
Pacchetto per la comunicazione con il database SQL.

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

func AddLibro(titolo string, autore, genere uint32) (uint32, error)
func AddGenere(nome string) (uint32, error)
func AddAutore(nome, cognome string) (uint32, error)
func AddPrestito(libro uint32, utente string data_prenotazione TipoData?, durata uint32) (uint32, error)

func SetHash(codice uint32, string hash) error
func SetRestituito(prestito uint32, data_restituzione TipoData?) error
```

