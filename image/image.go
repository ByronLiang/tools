package image

import (
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
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
