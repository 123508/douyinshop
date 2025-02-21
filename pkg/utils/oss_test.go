package utils

import "testing"

func TestUploadImages(t *testing.T) {
	localFilePath := "img.png"
	serviceName := "test"
	id := uint32(1)
	url, err := UploadImages(localFilePath, serviceName, id)
	if err != nil {
		t.Errorf("UploadImages() error = %v", err)
		return
	}
	t.Logf("UploadImages() url = %v", url)
}
