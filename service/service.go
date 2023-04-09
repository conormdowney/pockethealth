package service

import (
	"conordowney/pockethealth/dicom_wrapper"
	"conordowney/pockethealth/utils"
	"mime/multipart"
	"sync"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
)

// Upload executes a go routine per file to do
// the upload, tag retrieval and png conversion
func Upload(files []*multipart.FileHeader, tags []string) (map[string]map[string]*dicom.Element, error) {
	var wg sync.WaitGroup
	// a map to hold the result of each go routine by the name of the file
	elemFileNameMap := make(map[string]map[string]*dicom.Element)
	for i, fileHeader := range files {
		wg.Add(1)
		// add the name of the file so the map has the files in the order they were added by the user
		elemFileNameMap[fileHeader.Filename] = nil
		file, err := files[i].Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()
			elem, err := uploadFile(file, fileHeader, tags)
			if err != nil {
				return
			}
			// store the result keyed to the file name
			elemFileNameMap[fileHeader.Filename] = elem
		}(fileHeader)
	}

	wg.Wait()
	return elemFileNameMap, nil
}

// uploadFile handles the saving of the file being uploaded. Calls parseFile if a valid
// tag was supplied in the query request
func uploadFile(file multipart.File, fileHeader *multipart.FileHeader, tags []string) (map[string]*dicom.Element, error) {
	// use a uuid to name the file
	uuid := uuid.New()

	fileName, err := utils.UploadFile(file, fileHeader, uuid)
	if err != nil {
		return nil, err
	}
	// Get the data from the dcm file
	dataset, err := dicom_wrapper.ParseFile(fileName)
	if err != nil {
		return nil, err
	}
	// If there are any tags get the data associated to them tags
	var elem map[string]*dicom.Element
	if len(tags) > 0 {
		elem, err = dicom_wrapper.GetTags(tags, dataset)
		if err != nil {
			return nil, err
		}
	}
	// Convert the image from the dicom to a png
	err = dicom_wrapper.ConvertToPng(*dataset, uuid)
	if err != nil {
		return nil, err
	}

	return elem, nil
}
