package tiny_link

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"hash/crc32"
	"io"
	"strings"
)

var charts = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b",
	"c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
	"o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L",
	"M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z",
}

type TinyLink struct {
	RandomKey  string
	RangeIndex uint32
	Size       int // 随机码长度
	offset     int // 字节偏移量
}

func NewTinyLink(randomKey string, rangeIndex uint32, size int) *TinyLink {
	t := &TinyLink{
		RandomKey: randomKey,
		Size:      size,
	}
	if size > 31 {
		t.offset = 1
	} else {
		t.offset = 31 / size
	}
	if rangeIndex > uint32(len(charts)-1) {
		t.RangeIndex = uint32(len(charts) - 1)
	} else {
		t.RangeIndex = rangeIndex
	}
	return t
}

func (t *TinyLink) GenBigEndianLink(link string) ([4]string, error) {
	h := md5.New()
	h.Write([]byte(t.RandomKey + link))
	sMD5EncryptResult := h.Sum(nil)
	r := bytes.NewReader(sMD5EncryptResult)
	return t.GenTinyLink(r, binary.BigEndian)
}

func (t *TinyLink) GenLittleEndianLink(link string) ([4]string, error) {
	h := md5.New()
	h.Write([]byte(t.RandomKey + link))
	sMD5EncryptResult := h.Sum(nil)
	r := bytes.NewReader(sMD5EncryptResult)
	return t.GenTinyLink(r, binary.LittleEndian)
}

func (t *TinyLink) GenTinyLink(r io.Reader, order binary.ByteOrder) ([4]string, error) {
	var res = [4]string{}
	for i := 0; i < 4; i++ {
		var data uint32
		err := binary.Read(r, order, &data)
		if err != nil {
			return res, err
		}
		// 避免超出31位
		lHexLong := 0x3FFFFFFF & data
		chartGroup := make([]string, 0, t.Size)
		for j := 0; j < t.Size; j++ {
			index := t.RangeIndex & lHexLong
			chartGroup = append(chartGroup, charts[index])
			lHexLong = lHexLong >> t.offset
		}
		res[i] = strings.Join(chartGroup, "")
	}
	return res, nil
}

func HashCodeLink(link string) string {
	hashCode := crc32.ChecksumIEEE([]byte(link))
	// base62 生成随机码
	return GenBase62Code(int64(hashCode))
}
