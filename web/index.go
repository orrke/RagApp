package web

import (
	"net/http"
	"path"
	"text/template"
)

type IndexParams struct {
	Answer string
}

func ServeRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := IndexParams{
		Answer: "La réponse du modèle apparaîtra ici.",
	}

	tmpl, err := template.ParseFiles(path.Join("web", "templates", "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tmpl == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
