package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"time"
)

func GetImage(filename string) ([]byte, error) {
	fileByte, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileByte, nil
}

func OutImage() {
	file, err := os.Open("input.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	outfile, err := os.Create("output.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	if err := jpeg.Encode(outfile, img, nil); err != nil {
		log.Fatal(err)
	}
}

func DownloadImageHandle(w http.ResponseWriter, r *http.Request, fileName string, content []byte) {
	fm := mime.FormatMediaType("attachment", map[string]string{"filename": fileName})
	w.Header().Set("Content-Disposition", fm)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(content))
}
