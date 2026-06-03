package web

import "net/http"

// SetupRoutes is the function used to set up the routes for the web server.
func SetupRoutes() {
	//serve the assets publicly
	http.Handle("/assets/", http.StripPrefix("/assets/",
		http.FileServer(http.Dir("./web/assets")),
	))
	http.HandleFunc("/", ServeRoot)
}
