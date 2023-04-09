package handler

import (
	"bytes"
	"conordowney/pockethealth/router"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleUpload(t *testing.T) {
	os.Setenv("TESTING", "true")
	defer func() {
		os.Setenv("TESTING", "false")
		os.Remove("C:\\temp\\pngs\\tests")
		os.Remove("C:\\temp\\uploads\\tests")
	}()
	router := router.NewRouter()
	router.Post("/upload/", HandleUpload)

	filePath := "../test_files/1.dcm"
	fieldName := "files"
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	w, err := mw.CreateFormFile(fieldName, filePath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}

	// close the writer before making the request
	mw.Close()

	request, _ := http.NewRequest("POST", "/upload", body)
	request.Header.Add("content-type", mw.FormDataContentType())

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	err = request.ParseMultipartForm(200000)
	assert.Nil(t, err)
	formData := request.MultipartForm

	files := formData.File["files"]
	fmt.Print(len(files))

	assert.Equal(t, http.StatusOK, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	t.Log(response)
}

func TestHandleUploadWithoutMultipartForm(t *testing.T) {
	os.Setenv("TESTING", "true")
	defer func() {
		os.Setenv("TESTING", "false")
		os.RemoveAll("C:\\temp\\pngs\\tests")
		os.RemoveAll("C:\\temp\\uploads\\tests")
	}()
	router := router.NewRouter()
	router.Post("/upload/", HandleUpload)

	request, _ := http.NewRequest("POST", "/upload", nil)

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	t.Log(response)
}

func TestHandleUploadWithFilesInWrongAttribute(t *testing.T) {
	os.Setenv("TESTING", "true")
	defer func() {
		os.Setenv("TESTING", "false")
		os.Remove("C:\\temp\\pngs\\tests")
		os.Remove("C:\\temp\\uploads\\tests")
	}()
	router := router.NewRouter()
	router.Post("/upload/", HandleUpload)

	filePath := "../test_files/1.dcm"
	fieldName := "notfiles"
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	w, err := mw.CreateFormFile(fieldName, filePath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}

	// close the writer before making the request
	mw.Close()

	request, _ := http.NewRequest("POST", "/upload", body)
	request.Header.Add("content-type", mw.FormDataContentType())

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	t.Log(response)
}
