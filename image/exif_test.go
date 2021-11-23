package image

import (
	"bytes"
	"testing"

	"github.com/ByronLiang/tools/image/tiff"
)

func TestGetImage(t *testing.T) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		return
	}
	exifSize, sum := CalExifSize(fileByte)
	t.Log(exifSize, sum)

	_, startIndex, endIndex := GetExif(fileByte)
	t.Log("len: ", startIndex, endIndex)
	//t.Log(exifData)
}

func TestDecode(t *testing.T) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		return
	}
	exifData := GetExifData(fileByte, true)
	exifFile := bytes.NewReader(exifData)
	tiffGroup, err := tiff.Decode(exifFile)
	t.Log(tiffGroup.String())
	//tags, err := Decode(exifFile)
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	//tag, err := tags.Get(Orientation)
	//if err != nil || tag == nil {
	//	t.Fatal(err)
	//	return
	//}
	//t.Log(tag.String())
}

func TestGetDefineTag(t *testing.T) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		return
	}
	exifData := GetExifData(fileByte, true)
	exifFile := bytes.NewReader(exifData)
	// 获取指定 IFD 标签值
	tag, err := GetDefineTag(exifFile, Orientation)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("tag: ", tag.String())
}
