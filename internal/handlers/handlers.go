package handlers

import (
	"fmt"
	"net/http"
)

func HandleRootOr404(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	HandleHome(w, r)
}
func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "It works!")
}

func HandleLibri(w http.ResponseWriter, r *http.Request)        {}
func HandleAutori(w http.ResponseWriter, r *http.Request)       {}
func HandleGeneri(w http.ResponseWriter, r *http.Request)       {}
func HandleLogin(w http.ResponseWriter, r *http.Request)        {}
func HandleUtente(w http.ResponseWriter, r *http.Request)       {}
func HandlePrestito(w http.ResponseWriter, r *http.Request)     {}
func HandleRestituzione(w http.ResponseWriter, r *http.Request) {}
