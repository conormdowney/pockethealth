package handler

import (
	"conordowney/pockethealth/models"
	"conordowney/pockethealth/service"
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/gocraft/web"
	"github.com/suyashkumar/dicom"
)

// HandleUpload handles the case where upload is called.
func HandleUpload(rw web.ResponseWriter, req *web.Request) {
	// parses the body of a multi part form
	err := req.ParseMultipartForm(200000)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	formData := req.MultipartForm
	// get any files in the request
	files := formData.File["files"]
	if len(files) == 0 {
		http.Error(rw, "Bad Request: No files uploaded", http.StatusBadRequest)
		return
	}
	// validate files being uploaded are dcm files
	for _, file := range files {
		if filepath.Ext(file.Filename) != ".dcm" {
			http.Error(rw, "Bad Request: only dcm files accepted", http.StatusBadRequest)
			return
		}
	}
	// pass in the fiels and the tags to search for to the upload service
	elems, err := service.Upload(files, req.URL.Query()["tag"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(elems, rw)
}

// sendResponse handles sending the services response back to the client
func sendResponse(elems map[string]map[string]*dicom.Element, rw web.ResponseWriter) {
	rsp := models.Response{
		Result:     elems,
		StatusCode: http.StatusOK,
	}

	err := json.NewEncoder(rw).Encode(rsp.Result)
	if err != nil {
		// This would normally be logged somewhere
		panic("Error sending response: " + err.Error())
	}
}
