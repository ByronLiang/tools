package image

import (
	"net/http"
	"testing"
)

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
