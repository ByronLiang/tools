package material

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

func GetMd5FromUrl(url string) (string, []string) {
	bt, _ := GetCreativeUrlForBinaryData(url)
	mm := MD5Byte(bt)
	md5checks := ChunkByte(bt, len(bt), 5000)
	return mm, md5checks
}

func ChunkByte(bt []byte, fileSize, fileChunk int) []string {
	buff := bytes.NewBuffer(bt)
	totalPartsNum := int(math.Ceil(float64(fileSize) / float64(fileChunk)))
	md5checks := make([]string, 0, totalPartsNum)
	for i := 0; i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(fileChunk), float64(fileSize-i*fileChunk)))
		partBuffer := make([]byte, partSize)
		// chunk byte get
		buff.Read(partBuffer)
		md5checks = append(md5checks, MD5Byte(partBuffer))
		base64.StdEncoding.EncodeToString(partBuffer)
	}
	return md5checks
}

func GetMd5FromFile(filename string) (string, []string) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		return "", nil
	}
	defer file.Close()
	client := md5.New()
	_, err := io.Copy(client, file)
	if err != nil {
		return "", nil
	}
	mm := fmt.Sprintf("%x", client.Sum(nil))
	fileInfo, _ := file.Stat()
	md5checks := ChunkByteFromFile(filename, int(fileInfo.Size()), 5000)
	return mm, md5checks
}

func GetCreativeUrlForBinaryData(creativeUrl string) ([]byte, error) {
	var (
		req    *http.Request
		resp   *http.Response
		result []byte
		err    error
	)
	for i := 0; i < 3; i++ {
		req, err = http.NewRequest(http.MethodGet, creativeUrl, nil)
		if err != nil {
			continue
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("异常的http状态码: %d", resp.StatusCode)
			continue
		}

		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		break
	}
	return result, err
}

// MD5Byte sum bytes for md5
func MD5Byte(data []byte) string {
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ChunkByteFromFile(filename string, fileSize, fileChunk int) []string {
	file, openErr := os.Open(filename)
	if openErr != nil {
		return nil
	}
	defer file.Close()
	totalPartsNum := int(math.Ceil(float64(fileSize) / float64(fileChunk)))
	md5checks := make([]string, 0, totalPartsNum)
	for i := 0; i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(fileChunk), float64(fileSize-i*fileChunk)))
		partBuffer := make([]byte, partSize)
		// chunk byte get
		file.Read(partBuffer)
		md5checks = append(md5checks, MD5Byte(partBuffer))
		base64.StdEncoding.EncodeToString(partBuffer)
	}
	return md5checks
}
