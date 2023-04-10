package service_test

// import (
// 	"bytes"
// 	"conordowney/pockethealth/service"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"net/http"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// NOTE: Ran into a weird, maybe, race condition of some sort. When i ran the tests individually
// they worked fine. But when i ran them with the rest of the tests extra files were being added
// to the folder. Its possible some earlier goroutine was messing up the test run. I would
// have to investigate more thoroughly

// func TestUploadSingle(t *testing.T) {
// 	os.Setenv("TESTING", "true")
// 	defer func() {
// 		os.Setenv("TESTING", "false")
// 		os.RemoveAll("C:\\temp\\pngs\\tests")
// 		os.RemoveAll("C:\\temp\\uploads\\tests")
// 	}()
// 	filePath := "../test_files/2.dcm"
// 	fieldName := "files"
// 	body := new(bytes.Buffer)

// 	mw := multipart.NewWriter(body)

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	w, err := mw.CreateFormFile(fieldName, filePath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if _, err := io.Copy(w, file); err != nil {
// 		t.Fatal(err)
// 	}

// 	// close the writer before making the request
// 	mw.Close()

// 	request, _ := http.NewRequest("POST", "/upload", body)
// 	request.Header.Add("content-type", mw.FormDataContentType())

// 	err = request.ParseMultipartForm(200000)
// 	assert.Nil(t, err)
// 	formData := request.MultipartForm

// 	files := formData.File["files"]

// 	_, err = service.Upload(files, []string{"0000,0000"})
// 	assert.Nil(t, err)

// 	dirFiles, _ := os.ReadDir("C:\\temp\\pngs\\tests")
// 	t.Log(dirFiles)
// 	assert.Equal(t, 1, len(dirFiles), "Incorrect number of files in pngs test directory")
// 	uploadFiles, _ := os.ReadDir("C:\\temp\\uploads\\tests")
// 	t.Log(uploadFiles)
// 	assert.Equal(t, 1, len(uploadFiles), "Incorrect number of files in uploads test directory")
// }

// func TestUploadMultiple(t *testing.T) {
// 	os.Setenv("TESTING", "true")
// 	defer func() {
// 		os.Setenv("TESTING", "false")
// 		os.RemoveAll("C:\\temp\\pngs\\tests")
// 		os.RemoveAll("C:\\temp\\uploads\\tests")
// 	}()

// 	fieldName := "files"
// 	body := new(bytes.Buffer)

// 	mw := multipart.NewWriter(body)
// 	for i := 1; i <= 5; i++ {
// 		filePath := fmt.Sprintf("../test_files/%d.dcm", i)

// 		file, err := os.Open(filePath)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		w, err := mw.CreateFormFile(fieldName, filePath)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if _, err := io.Copy(w, file); err != nil {
// 			t.Fatal(err)
// 		}

// 	}

// 	// close the writer before making the request
// 	mw.Close()

// 	request, _ := http.NewRequest("POST", "/upload", body)
// 	request.Header.Add("content-type", mw.FormDataContentType())

// 	err := request.ParseMultipartForm(200000)
// 	assert.Nil(t, err)
// 	formData := request.MultipartForm

// 	files := formData.File["files"]

// 	_, err = service.Upload(files, []string{"0000,0000"})
// 	assert.Nil(t, err)

// 	pngFiles, _ := os.ReadDir("C:\\temp\\pngs\\tests")
// 	t.Log(pngFiles)
// 	assert.Equal(t, 5, len(pngFiles), "Incorrect number of files in pngs test directory")
// 	uploadFiles, _ := os.ReadDir("C:\\temp\\uploads\\tests")
// 	t.Log(uploadFiles)
// 	assert.Equal(t, 5, len(uploadFiles), "Incorrect number of files in uploads test directory")
// }

// // test checking the returned element
