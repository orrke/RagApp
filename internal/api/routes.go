package api

import "net/http"

func SetupRoutes(apiPrefix string) {
	http.HandleFunc(apiPrefix+"/search", SearchDocument)
	http.HandleFunc(apiPrefix+"/search/raw", SearchRawDocument)
	http.HandleFunc(apiPrefix+"/update", UpdateIndex)
	http.HandleFunc(apiPrefix+"/reset", ResetIndex)
	http.HandleFunc(apiPrefix+"/reconfig", ReloadConfig)
}
