package dicom_wrapper_test

import (
	"conordowney/pockethealth/dicom_wrapper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	dataset, err := dicom_wrapper.ParseFile("../test_files/1.dcm")
	assert.Nil(t, err)
	assert.NotNil(t, dataset)

	_, err = dicom_wrapper.ParseFile("../test_files/doesnotexist.dcm")
	assert.NotNil(t, err)
}

func TestGetTag(t *testing.T) {
	dataset, err := dicom_wrapper.ParseFile("../test_files/1.dcm")
	assert.Nil(t, err)

	elem, err := dicom_wrapper.GetTags([]string{"7fe0,0010"}, dataset)
	assert.Nil(t, err)
	assert.NotNil(t, elem)

	compareTag, err := dicom_wrapper.CreateTag("7fe0,0010")
	assert.Nil(t, err)
	compare := elem["7fe0,0010"].Tag.Equals(*compareTag)
	assert.True(t, compare)
}

func TestGetTagErrors(t *testing.T) {
	dataset, err := dicom_wrapper.ParseFile("../test_files/1.dcm")
	assert.Nil(t, err)

	// looking for a tag that doesnt exist on the element skips that tag
	// bu produces no error
	_, err = dicom_wrapper.GetTags([]string{"ffff,ffff"}, dataset)
	assert.Nil(t, err)

	// passing in non hex values should give an error
	_, err = dicom_wrapper.GetTags([]string{"xzxz,xyxy"}, dataset)
	assert.NotNil(t, err)

	// passing in non hex values should give an error
	_, err = dicom_wrapper.GetTags([]string{"0000,xyxy"}, dataset)
	assert.NotNil(t, err)
}

func TestConvertToPng(t *testing.T) {
	dataset, err := dicom_wrapper.ParseFile("../test_files/1.dcm")
	assert.Nil(t, err)

	uuid := uuid.New()
	err = dicom_wrapper.ConvertToPng(*dataset, uuid)
	assert.Nil(t, err)
}
