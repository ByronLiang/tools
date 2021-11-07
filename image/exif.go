package image

import (
	"bytes"
	"encoding/binary"

	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
)

func RemoveExif(data []byte) ([]byte, []byte, error) {
	jmp := jpegstructure.NewJpegMediaParser()
	if jmp.LooksLikeFormat(data) {
		sl, err := jmp.ParseBytes(data)
		if err != nil {
			return data, nil, err
		}
		_, rawExif, err := sl.Exif()
		if err != nil {
			// 不存在图片exif信息
			return data, nil, nil
		}
		var filtered []byte
		startExifBytes := 0
		endExifBytes := 0
		if bytes.Contains(data, rawExif) {
			for i := 0; i < len(data)-len(rawExif); i++ {
				if bytes.Compare(data[i:i+len(rawExif)], rawExif) == 0 {
					startExifBytes = i
					endExifBytes = i + len(rawExif)
					break
				}
			}
			fill := make([]byte, len(data[startExifBytes:endExifBytes]))
			copy(data[startExifBytes:endExifBytes], fill)
		}
		filtered = data
		return filtered, rawExif, nil
	}
	return data, nil, nil
}

func GetExif(data []byte) ([]byte, int, int) {
	startExifBytesIndex := 0
	endExifBytesIndex := 0
	jmp := jpegstructure.NewJpegMediaParser()
	if jmp.LooksLikeFormat(data) {
		sl, err := jmp.ParseBytes(data)
		if err != nil {
			return nil, startExifBytesIndex, endExifBytesIndex
		}
		_, rawExif, err := sl.Exif()
		if err != nil {
			// 不存在图片exif信息
			return nil, startExifBytesIndex, endExifBytesIndex
		}
		if bytes.Contains(data, rawExif) {
			for i := 0; i < len(data)-len(rawExif); i++ {
				if bytes.Compare(data[i:i+len(rawExif)], rawExif) == 0 {
					startExifBytesIndex = i
					endExifBytesIndex = i + len(rawExif)
					break
				}
			}
		}
		return rawExif, startExifBytesIndex, endExifBytesIndex
	}
	return nil, startExifBytesIndex, endExifBytesIndex
}

func bytesToInt16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func CalExifSize(data []byte) (int, int) {
	exifSize := 0
	sum := 0
	if data[2] == 0xFF {
		switch data[3] {
		case 0xE0:
			jfifSize := bytesToInt16(data[4:6])
			if data[jfifSize+4] == 0xFF && data[jfifSize+5] == 0xE1 {
				offset := int(jfifSize) + 4
				exifSize = int(bytesToInt16(data[offset+2 : offset+4]))
				sum = offset + exifSize + 2
			}
			break
		case 0xE1:
			offset := 2
			exifSize = int(bytesToInt16(data[offset+2 : offset+4]))
			sum = offset + exifSize + 2
		}
	}
	return exifSize, sum
}
