package main

import (
	"net/http"
)

func (app *application) trackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}