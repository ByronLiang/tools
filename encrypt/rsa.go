package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"os"
	"strings"
	"syscall"
)

const (
	privateFileName  = "private.pem"
	publicFileName   = "public.pem"
	privateKeyPrefix = "RSA PRIVATE KEY "
	publicKeyPrefix  = "RSA PUBLIC KEY "
)

const DefaultKeySize = 2048

var RsaEncrypt *rsaEncrypt

type rsaEncrypt struct {
	publicKey  []byte
	privateKey []byte
	limitSize int
}

// 读取public.pem, private.pem密匙文件
func NewRsaEncrypt(publicKey, privateKey string) {
	RsaEncrypt = &rsaEncrypt{
		publicKey:  []byte(publicKey),
		privateKey: []byte(privateKey),
	}
}

func InitRsaEncrypt() {
	RsaEncrypt = &rsaEncrypt{}
}

func GetRsaKey(keySize int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}
	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create(privateFileName)
	if err != nil {
		return err
	}
	defer privateFile.Close()
	privateBlock := pem.Block{
		Type:  privateKeyPrefix,
		Bytes: x509PrivateKey,
	}

	if err = pem.Encode(privateFile, &privateBlock); err != nil {
		return err
	}
	publicKey := privateKey.PublicKey
	x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	publicFile, _ := os.Create(publicFileName)
	defer publicFile.Close()
	publicBlock := pem.Block{
		Type:  publicKeyPrefix,
		Bytes: x509PublicKey,
	}
	if err = pem.Encode(publicFile, &publicBlock); err != nil {
		return err
	}
	return nil
}

func (r *rsaEncrypt) encrypt(content string) (string, error) {
	plainText := []byte(content)
	block, _ := pem.Decode(r.publicKey)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipherText), nil
}

func (r *rsaEncrypt) decrypt(cryptContent string) (string, error) {
	// 私匙解密
	block, _ := pem.Decode(r.privateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	cryptText, err := hex.DecodeString(cryptContent)
	if err != nil {
		return "", err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cryptText)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

// 获取可加密明文的最大长度
func (r *rsaEncrypt)GetLimitMsgSize() (int, error) {
	block, _ := pem.Decode(r.publicKey)
	if block == nil {
		return 0, errors.New("public key decode fail")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return 0, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	size := publicKey.Size() - 12
	return size, nil
}

func (r *rsaEncrypt) SetLimitMsgSize() error {
	size, err := r.GetLimitMsgSize()
	if err != nil {
		return err
	}
	r.limitSize = size
	return nil
}

// 识别长文本并分组加密
func (r *rsaEncrypt) encryptLongText(content string) ([]string, error) {
	if r.limitSize == 0 {
		// 加密长度限制异常
		return nil, errors.New("limitSize error")
	}
	// 检测是否超出可加密字符长度
	if len(content) > r.limitSize {
		i := 1
		offset := len(content) / r.limitSize
		restSize := len(content) % r.limitSize
		encryptListSize := offset
		if restSize > 0 {
			// 预留一位存放剩余加密内容
			encryptListSize ++
		}
		contentEncryptList := make([]string, encryptListSize)
		for i = 1; i <= offset; i++ {
			start := (i-1) * r.limitSize
			end := i * r.limitSize
			c := make([]byte, r.limitSize)
			copy(c, content[start:end])
			dataE, err := r.encrypt(string(c))
			if err != nil {
				contentEncryptList[i-1] = ""
				continue
			}
			contentEncryptList[i-1] = dataE
		}
		// 对剩余内容进行加密
		if restSize > 0 {
			c := make([]byte, restSize)
			copy(c, content[(i - 1) * r.limitSize:])
			dataE, err := r.encrypt(string(c))
			if err != nil {
				contentEncryptList[i - 1] = ""
			} else {
				contentEncryptList[i - 1] = dataE
			}
		}
		return contentEncryptList, nil
	}
	// 无需分组加密
	dataE, err := RsaEncrypt.encrypt(content)
	if err != nil {
		return nil, err
	}
	return []string{dataE}, nil
}

func (r *rsaEncrypt) decryptLongText(cryptContentList []string) string {
	contentDecryptList := make([]string, len(cryptContentList))
	for j := 0; j < len(cryptContentList); j++ {
		if cryptContentList[j] != "" {
			c, err := r.decrypt(cryptContentList[j])
			if err == nil {
				contentDecryptList[j] = c
			}
		}
	}
	return strings.Join(contentDecryptList, "")
}

// 不建议生产环境使用，适用测试环境
func (r *rsaEncrypt) GetKeyFromFile(path, filename string) error {
	file := path + filename
	d, err := syscall.Open(file, syscall.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(d)
	var buf [2048]byte
	n, err := syscall.Read(d, buf[:])
	if err != nil {
		return err
	}
	keyBuf := make([]byte, n)
	if filename == privateFileName {
		copy(keyBuf, buf[:n])
		r.privateKey = keyBuf
	}
	if filename == publicFileName {
		copy(keyBuf, buf[:n])
		r.publicKey = keyBuf
	}
	return nil
}


