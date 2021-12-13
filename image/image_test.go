package image

import (
	"bytes"
	"net/http"
	"testing"
)

func TestOutImage(t *testing.T) {
	OutImage("./sample/demo.jpg", "./sample/output.jpg")
}

func TestGetImageInfo(t *testing.T) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	fileReader := bytes.NewReader(fileByte)
	w, h, format, err := GetImageInfo(fileReader)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("image w: %d; h: %d; format: %s", w, h, format)
}

func TestDownloadImageHandle(t *testing.T) {
	http.HandleFunc("/download", downloadImageHandle)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		t.Error(err)
	}
}

func downloadImageHandle(w http.ResponseWriter, r *http.Request) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// RemoveExif(fileByte)
	RemoveExifSkipOrientation(fileByte)
	DownloadImageHandle(w, r, "image.jpg", fileByte)
}
