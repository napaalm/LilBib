# template
Qui sono organizzati i template Go/HTML per le pagine servite dal server.
Vedere questa [pagina](https://golang.org/pkg/html/template/) per l'utilizzo.

## Campi disponibili nei template
#### Pagina `/libro`
```go
type Libro struct {
	Codice    uint32
	Titolo    string
	Autore    uint32
	Genere    uint32
	Prenotato bool
	Hash      string
}
```

#### Pagina `/libri`
```go
struct {
	Pagina uint64
	Libri  []Libro
}
```

#### Pagina `/autori`
```go
struct {
	Iniziale byte
	Autori   []Autore
}

type Autore struct {
	Codice  uint32
	Nome    string
	Cognome string
}
```

#### Pagina `/generi`
```go
struct {
	Generi []Genere
}
```
