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
}

func TestDecode(t *testing.T) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		return
	}
	exifData, existExif := GetExifData(fileByte, true)
	if existExif {
		exifFile := bytes.NewReader(exifData)
		tiffGroup, err := tiff.Decode(exifFile)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(tiffGroup.String())
	}
}

func TestGetDefineTag(t *testing.T) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		return
	}
	exifData, isExistExif := GetExifData(fileByte, true)
	if !isExistExif {
		return
	}
	exifFile := bytes.NewReader(exifData)
	// 获取指定 IFD 标签值
	tag, err := GetDefineTag(exifFile, Orientation)
	if err != nil {
		t.Error(err)
		return
	}
	if tag != nil {
		t.Log("val: ", tag.String())
	}
}

func TestGetImageOrientation(t *testing.T) {
	fileByte, err := GetImage("./sample/demo.jpg")
	if err != nil {
		return
	}
	tagVal, err := GetImageOrientation(fileByte)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("tagVal: ", tagVal)
}
