package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

var AesEncrypt *aesEncrypt

type aesEncrypt struct {
	name string
	key  []byte
}

func NewAesEncrypt(name, key string) {
	// 校验key 长度必须 16/24/32
	AesEncrypt = &aesEncrypt{name: name, key: []byte(key)}
}

func (a *aesEncrypt) GetName() string {
	return a.name
}

func (a *aesEncrypt) GetKey() string {
	return string(a.key)
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	return append(ciphertext, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	return origData[:(length - int(origData[length-1]))]
}

func paddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}

func (a *aesEncrypt) Decrypt(text string) (string, error) {
	// hex解析
	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	//和js的key补码方法一致
	pkey := paddingLeft(a.key, '0', 16)
	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return "", err
	}
	blockModel := cipher.NewCBCDecrypter(block, pkey) //和前端代码对应:   mode: CryptoJS.mode.CBC,// CBC算法
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = pkcs5UnPadding(plantText) //和前端代码对应:  padding: CryptoJS.pad.Pkcs7
	return string(plantText), nil
}

func (a *aesEncrypt) Encrypt(raw string) (string, error) {
	origData := []byte(raw)
	k := paddingLeft(a.key, '0', 16)
	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pkcs5Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k)
	// 创建数组
	crypts := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(crypts, origData)
	// 返回hex
	return hex.EncodeToString(crypts), nil
}

