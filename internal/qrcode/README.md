# qrcode
Pacchetto per la generazione dei QR code destinati ai libri.

## Funzionamento

Le funzioni in questo package generano le password per uno o più libri e le salvano un in QR code. Il formato delle immagini è png.

## Tipi
```go
type QRCodeLibro
```

## Funzioni
```go
func CreateQRCode(id uint32) (QRCodeLibro, error) // genera e ritorna il QR code del libro con l'id passato in argomento
func GeneratePage(ids []uint32) (string, error) // ritorna una pagina HTML con i codici QR
```

