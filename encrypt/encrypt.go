package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func Pkcs5Padding(src []byte, size int) (dst []byte) {
	rbs := size - len(src)%size
	return append(src, bytes.Repeat([]byte{byte(rbs)}, rbs)...)
}

func Pkcs5UnPadding(src []byte) (dst []byte, err error) {
	if len(src) == 0 {
		return nil, fmt.Errorf("empty PKCS5Padding")
	}
	paddingNum := int(src[len(src)-1])
	if paddingNum > len(src) {
		return nil, fmt.Errorf("invalid PKCS5Padding: %v", src)
	}

	return src[:len(src)-paddingNum], nil
}

func base64UrlSafeEncode(source []byte) string {
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}

// base64 safe url decode
func base64UrlSafeDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	return base64.URLEncoding.DecodeString(data)
}

func AesEBCEncrypt(src string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	if src == "" {
		return "", errors.New("plain content empty")
	}
	srcBytes := []byte(src)
	srcBytes = Pkcs5Padding(srcBytes, block.BlockSize())
	dst := make([]byte, len(srcBytes))
	cryptBlocksErr := cryptBlocks(dst, srcBytes, block)
	if cryptBlocksErr != nil {
		return "", cryptBlocksErr
	}
	return base64UrlSafeEncode(dst), nil
}

func AesEBCDecrypt(encPrice string, key string) (string, error) {
	priceByte, err := base64UrlSafeDecode(encPrice)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	origData := make([]byte, len(priceByte))
	decryptBlocksErr := decryptBlocks(origData, priceByte, block)
	if decryptBlocksErr != nil {
		return "", decryptBlocksErr
	}
	origData, err = Pkcs5UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}

func cryptBlocks(dst, srcBytes []byte, block cipher.Block) error {
	if len(srcBytes)%block.BlockSize() != 0 {
		return errors.New("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(srcBytes) {
		return errors.New("crypto/cipher: output smaller than input")
	}
	for len(srcBytes) > 0 {
		block.Encrypt(dst, srcBytes[:block.BlockSize()])
		srcBytes = srcBytes[block.BlockSize():]
		dst = dst[block.BlockSize():]
	}
	return nil
}

func decryptBlocks(dst, src []byte, b cipher.Block) error {
	if len(src)%b.BlockSize() != 0 {
		return errors.New("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		return errors.New("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		b.Decrypt(dst, src[:b.BlockSize()])
		src = src[b.BlockSize():]
		dst = dst[b.BlockSize():]
	}
	return nil
}
