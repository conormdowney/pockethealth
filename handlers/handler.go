package handler

import (
	"conordowney/pockethealth/models"
	"conordowney/pockethealth/service"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gocraft/web"
)

// HandleUpload handles the case where upload is called.
func HandleUpload(rw web.ResponseWriter, req *web.Request) {
	files, err := validateRequest(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// pass in the files and the tags to search for to the upload service
	elems, err := service.Upload(files, req.URL.Query()["tag"])
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(elems, rw)
}

// HandleConvertToPng handles the case where upload is called.
func HandleConvertToPng(rw web.ResponseWriter, req *web.Request) {

	files, err := validateRequest(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	pngs, err := service.ConvertToPng(files)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(pngs, rw)
}

// sendResponse handles sending the services response back to the client
func sendResponse(body interface{}, rw web.ResponseWriter) {
	rsp := models.Response{
		Result:     body,
		StatusCode: http.StatusOK,
	}

	err := json.NewEncoder(rw).Encode(rsp.Result)
	if err != nil {
		// This would normally be logged somewhere
		panic("Error sending response: " + err.Error())
	}
}

// validateRequest validates the request is valid. Also extracts the files from the request
func validateRequest(req *web.Request) ([]*multipart.FileHeader, error) {
	// parses the body of a multi part form
	err := req.ParseMultipartForm(200000)
	if err != nil {
		return nil, err
	}
	formData := req.MultipartForm
	// get any files in the request
	files := formData.File["files"]
	if len(files) == 0 {
		return nil, errors.New("Bad Request: No files uploaded")
	}
	// validate files being uploaded are dcm files
	for _, file := range files {
		if filepath.Ext(file.Filename) != ".dcm" {
			return nil, errors.New("Bad Request: only dcm files accepted")
		}
	}

	return files, nil
}
