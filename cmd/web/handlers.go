package main

import (
	"errors"
	"fmt"
	"github.com/Reticent93/snips/pkg/models"
	//"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snips.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snips: s,
	})

	//data := &templateData{Snips: s}
	//
	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}

	//Parsing the html templates
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//http.Error(w, "Internal Server Error", 500)
	//return
	//}

	//Execute method writes content as the response body
	//err = ts.Execute(w, data)
	//if err != nil {
	//	app.serverError(w, err)
	//}
}

func (app *application) showSnip(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snips.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snip: s,
	})

}

func (app *application) createSnipForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snip..."))
}

func (app *application) createSnip(w http.ResponseWriter, r *http.Request) {

	title := "O yokai"
	content := "O wanka\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "14 days"

	id, err := app.snips.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snip/%d", id), http.StatusSeeOther)
}
