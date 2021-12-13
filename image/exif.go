package image

import (
	"bytes"
	"encoding/binary"
)

func CalExifSize(data []byte) (int, int) {
	exifSize := 0
	sum := 0
	if data[2] == 0xFF {
		switch data[3] {
		case 0xE0:
			jfifSize := bytesToInt16(data[4:6])
			if data[jfifSize+4] == 0xFF && data[jfifSize+5] == 0xE1 {
				offset := int(jfifSize) + 4
				exifSize = int(bytesToInt16(data[offset+2:offset+4])) - 2
				sum = offset + exifSize + 4
			}
			break
		case 0xE1:
			offset := 2
			exifSize = int(bytesToInt16(data[offset+2:offset+4])) - 2
			sum = offset + exifSize + 4
		}
	}
	return exifSize, sum
}

func bytesToInt16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func GetExifData(data []byte, isOnlyContent bool) ([]byte, bool) {
	exifSize, sumSize := CalExifSize(data)
	if exifSize == 0 || sumSize == 0 {
		return nil, false
	}
	startIndex := sumSize - exifSize
	// 只提取exif 内容
	if isOnlyContent {
		startIndex += 6
		exifSize -= 6
	}
	// 从源数据拷贝exif数据
	exifData := make([]byte, exifSize)
	copy(exifData, data[startIndex:sumSize])
	return exifData, true
}

func RemoveExif(data []byte) {
	exifSize, sumSize := CalExifSize(data)
	if exifSize == 0 || sumSize == 0 {
		return
	}
	startIndex := sumSize - exifSize + 6
	// 覆盖exif 数据
	fill := make([]byte, len(data[startIndex:sumSize]))
	copy(data[startIndex:sumSize], fill)
	return
}

func RemoveExifSkipOrientation(data []byte) {
	exifSize, sumSize := CalExifSize(data)
	if exifSize == 0 || sumSize == 0 {
		return
	}
	// EXIF IFD 标签起始字节下标
	startIndex := sumSize - exifSize + 6
	// 覆盖exif 数据
	fill := make([]byte, len(data[startIndex:sumSize]))
	// 提前前三个 EXIF 的 IDF 字节数据 前10个字节为 EXIF 标签基本信息 每12个字节存储一个 IDF 标签信息
	copy(fill[:46], data[startIndex:startIndex+46])
	// 设置 EXIF 的 IDF 数目值修改为1
	//tagTotalIdf := data[startIndex+8 : startIndex+10]
	//binary.BigEndian.PutUint16(tagTotalIdf, uint16(1))
	// 覆盖源字节数据
	//copy(fill[:10], data[startIndex:startIndex+10])
	// 单独提取 图片旋转角度
	//copy(fill[10:22], data[startIndex+34:startIndex+46])
	copy(data[startIndex:sumSize], fill)
	return
}

func GetImageOrientation(fileByte []byte) (int, error) {
	exifData, isExistExif := GetExifData(fileByte, true)
	if !isExistExif {
		return 0, nil
	}
	exifFile := bytes.NewReader(exifData)
	// 获取指定 IFD 标签值
	tag, err := GetDefineTag(exifFile, Orientation)
	if err != nil {
		return 0, err
	}
	if tag != nil {
		tagVal, err := tag.Int(0)
		if err != nil {
			return 0, err
		}
		return tagVal, nil
	}
	return 0, nil
}
