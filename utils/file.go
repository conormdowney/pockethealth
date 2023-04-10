package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// UploadFile uploads a file and returns the file name
func UploadFile(file multipart.File, fileHeader *multipart.FileHeader, uuid uuid.UUID) (string, error) {
	// During unit tests the images are stored in a tests folder inside the uploads folder.
	// This is to allow the tests to remove the images when they are cleaning up
	uploadDir := "C:/temp/uploads"
	if os.Getenv("TESTING") == "true" {
		uploadDir += "/tests"
	}
	// Create the uploads folder if it doesn't already exist
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	// get the filename in the form uuid.file extension of the uploaded file e.g. 1234-5678.dcm
	fileName := fmt.Sprintf("%s/%s%s", uploadDir, uuid, filepath.Ext(fileHeader.Filename))
	// Create a new file in the uploads directory
	dst, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	// Copy the uploaded file to the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
