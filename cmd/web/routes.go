package main

import (
	"github.com/go-chi/chi"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := chi.NewRouter()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snip/create", http.HandlerFunc(app.createSnipForm))
	mux.Post("/snip/create", http.HandlerFunc(app.createSnip))
	mux.Get("/snip/:id", http.HandlerFunc(app.showSnip))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
