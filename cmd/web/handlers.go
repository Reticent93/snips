package main

import (
	"errors"
	"fmt"
	"github.com/Reticent93/snips/pkg/models"
	"html/template"

	//"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snips.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snips: s}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	//Parsing the html templates
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		//http.Error(w, "Internal Server Error", 500)
		return
	}

	//Execute method writes content as the response body
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnip(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	data := &templateData{Snip: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	//Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snip: s,
	})

}

func (app *application) createSnip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set(
			"Allow",
			http.MethodPost,
		)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O yokai"
	content := "O wanka\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "14 days"

	id, err := app.snips.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snip?id=%d", id), http.StatusSeeOther)
}
