package main

import (
	handler "conordowney/pockethealth/handlers"
	"conordowney/pockethealth/router"
	"net/http"
)

type Context struct{}

func main() {
	router := router.NewRouter()

	// requests can be sent in the form: upload?tag=0400,0565&tag=0010,0020
	// and can include multiple files and tags
	router.Post("/upload/", handler.HandleUpload)

	http.ListenAndServe("localhost:3000", router)
}
