package web

import "net/http"

func SetupRoutes() {
	http.Handle("/assets/", http.StripPrefix("/assets/",
		http.FileServer(http.Dir("./web/assets")),
	))
	http.HandleFunc("/", ServeRoot)
}
