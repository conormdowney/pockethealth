package dicom_wrapper

import (
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// ParseFile parses a file to allow searching for a given tag
func ParseFile(fileName string) (*dicom.Dataset, error) {
	dataset, err := dicom.ParseFile(fileName, nil)
	if err != nil {
		return nil, err
	}

	//fmt.Println(dataset)

	return &dataset, nil
}

// GetTags returns data for a list of tags if they exist on the dataset. If they
// don't they are just ignored in the return map.
// The returned object is a map of tag code -> tag data e.g.
// "{0010,0020":{"tag":{"Group":16,...}}, "0400,0565":{"tag":{"Group":1024,...}}}
func GetTags(queryTags []string, dataset *dicom.Dataset) (map[string]*dicom.Element, error) {
	tagMap := make(map[string]*dicom.Element)
	for _, splitTag := range queryTags {
		// create a tag for use in the FindElementByTagNested function
		tag, err := CreateTag(splitTag)
		if err != nil {
			return nil, err
		}
		// Get the data from the dataset for the given tag
		elem, err := dataset.FindElementByTagNested(*tag)
		if err != nil {
			// if the error is that the tag could not be found, just continue on to the next item
			if strings.Contains(err.Error(), "element not found") {
				continue
			} else {
				return nil, err
			}
		}

		tagMap[splitTag] = elem
	}

	return tagMap, nil
}

// ConvertToPng converts a dicom image to png nd stores it locally
func ConvertToPng(dataset dicom.Dataset, uuid uuid.UUID) error {
	pngDir := "C:\\temp\\pngs"
	// When running tests, store the items in a tests folder that should be
	// deleted after each test
	if os.Getenv("TESTING") == "true" {
		pngDir += "\\tests"
	}
	err := os.MkdirAll(pngDir, os.ModePerm)
	if err != nil {
		return err
	}
	// Find the data related to images
	pixelDataElement, _ := dataset.FindElementByTag(tag.PixelData)
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for i, fr := range pixelDataInfo.Frames {
		img, err := fr.GetImage()
		if err != nil {
			return err
		}
		// The files used in this assignment only have single frames, but in the case where
		// there are multiple the files will be stored as <uuid>_0.
		// Ideally this path, or in a cloud setting the S3 (or whatever Azures equivalent is)
		// would be stored in a database in an entry for the image
		f, err := os.Create(fmt.Sprintf("%s\\%s_%d.png", pngDir, uuid, i))
		if err != nil {
			return err
		}
		err = png.Encode(f, img)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTag creates a tag from the codes passed in as query parameters
func CreateTag(tagString string) (*tag.Tag, error) {
	tagParts := strings.Split(tagString, ",")
	var base = 16
	var size = 16
	// tags are strings and need to be uint16's for creating the tag
	val1, err := strconv.ParseUint(tagParts[0], base, size)
	if err != nil {
		return nil, err
	}
	val2, err := strconv.ParseUint(tagParts[1], base, size)
	if err != nil {
		return nil, err
	}

	newTag := tag.Tag{Group: uint16(val1), Element: uint16(val2)}
	return &newTag, nil
}
