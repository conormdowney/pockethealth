package utils_test

import (
	"bytes"
	"conordowney/pockethealth/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	os.Setenv("TESTING", "true")
	defer func() {
		os.Setenv("TESTING", "false")
		os.RemoveAll("C:\\temp\\uploads\\tests")
	}()
	filePath := "../test_files/2.dcm"
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

	err = request.ParseMultipartForm(200000)
	assert.Nil(t, err)
	formData := request.MultipartForm

	files := formData.File["files"]
	fileToUpload, err := files[0].Open()
	assert.Nil(t, err)

	uuid := uuid.New()
	fileName, err := utils.UploadFile(fileToUpload, files[0], uuid)
	assert.Nil(t, err)

	_, err = os.Stat(fileName)
	assert.Nil(t, err)
}
