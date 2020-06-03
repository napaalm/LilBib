# db
Pacchetto per la comunicazione con il database SQL.
Utilizzare le seguenti librerie: [database/sql](https://git.antonionapolitano.eu/napaalm/LilBib/src/master) e [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).

## Tipi
```go
type Libro
type Autore
type Genere
type Prestito
```

## Funzioni
```go
func InizializzaDB() (error)
func ChiudiDB()

func GetLibro(codice uint32) (Libro, error)
func GetLibri(page uint32) ([]Libri, error)
func GetAutori(iniziale uint8) ([]Autore, error)
func GetGeneri() ([]Genere, error)
func GetPrestiti(utente string) ([]Prestito, error)
func GetAssegnatario(cod uint32) (string, error)

func RicercaLibri(nome string, autore, genere []uint32, page uint16) ([]Libro, error)
func RicercaAutori(nome string) ([]Autore, error)
func RicercaGeneri(nome string) ([]Genere, error)

func AddLibro(titolo string, autore, genere uint32) (uint32, error)
func AddGenere(nome string) (uint32, error)
func AddAutore(nome, cognome string) (uint32, error)
func AddPrestito(libro uint32, utente string, data_prenotazione time.Time, durata uint32) (uint32, error)

func SetHash(codice uint32, hash string) error
func SetRestituito(prestito uint32, data_restituzione time.Time) error

func RemoveLibro(codice uint32) error
func RemoveGenere(codice uint32) error
func RemoveAutore(codice uint32) error
func RemovePrestito(codice uint32) error

func LibriPrenotati() (uint32, error)
func LibriDisponibili() (uint32, error)
```

