package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

func GetImage(filename string) ([]byte, error) {
	fileByte, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileByte, nil
}

func GetImageInfo(imageReader *bytes.Reader) (int, int, string, error) {
	_, err := imageReader.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, "", err
	}
	c, format, err := image.DecodeConfig(imageReader)
	if err != nil {
		return 0, 0, "", err
	}
	return c.Width, c.Height, format, nil
}

func OutImage(filename, outputFileName string) error {
	fileByte, err := GetImage(filename)
	if err != nil {
		return err
	}
	orientation, err := GetImageOrientation(fileByte)
	if err != nil {
		return err
	}
	fileReader := bytes.NewReader(fileByte)
	img, _, err := image.Decode(fileReader)
	if err != nil {
		return err
	}
	if orientation == 3 {
		img = imaging.Rotate180(img)
	}
	if orientation == 6 {
		img = imaging.Rotate270(img)
	}
	if orientation == 8 {
		img = imaging.Rotate90(img)
	}
	outfile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outfile.Close()
	if err := jpeg.Encode(outfile, img, nil); err != nil {
		return err
	}
	return nil
}

// Http 请求下以流形式下载图片封装方法
func DownloadImageHandle(w http.ResponseWriter, r *http.Request, fileName string, content []byte) {
	fm := mime.FormatMediaType("attachment", map[string]string{"filename": fileName})
	w.Header().Set("Content-Disposition", fm)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(content))
}
