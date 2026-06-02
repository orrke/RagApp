package web

import (
	"net/http"
	"path"
	"text/template"
)

// ServeRoot is the function that responds to calls at the root. It returns an HTML page.
func ServeRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//get the template HTML file for the index
	tmpl, err := template.ParseFiles(path.Join("web", "templates", "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tmpl == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	//serve the template file
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
