package image

import "testing"

func TestGetImage(t *testing.T) {
	fileByte, err := GetImage("./sample/tick.JPG")
	if err != nil {
		return
	}
	exifSize, sum := CalExifSize(fileByte)
	t.Log(exifSize, sum)
	_, startIndex, endIndex := GetExif(fileByte)
	t.Log("len: ", startIndex, endIndex)
	//t.Log(exifData)
}
