package encrypt

import (
	"fmt"
	"syscall"
	"testing"
)

func TestGetRsaKey(t *testing.T) {
	err := GetRsaKey(DefaultKeySize)
	if err != nil {
		t.Error(err)
	}
}

func TestNewRsaEncrypt(t *testing.T) {
	initRsaEncryptFromFile()
	limitSize, err := RsaEncrypt.GetLimitMsgSize()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("limit size: ", limitSize)
	dataE, err := RsaEncrypt.encrypt("abc")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataE)
		c, err := RsaEncrypt.decrypt(dataE)
		if err == nil {
			t.Log(c)
		}
	}
}

func initRsaEncryptFromFile()  {
	InitRsaEncrypt()
	// 从文件读私匙
	err := RsaEncrypt.GetKeyFromFile("./", privateFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 从文件读公匙
	err = RsaEncrypt.GetKeyFromFile("./", publicFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestAesEncrypt_Encrypt(t *testing.T) {
	initRsaEncryptFromFile()
	err:= RsaEncrypt.SetLimitMsgSize()
	if err != nil {
		t.Error(err)
		return
	}
	// 长文本加密
	txt, err := readFile("./", "longText.txt")
	if err != nil {
		t.Error(err)
		return
	}
	contentEncryptList, err := RsaEncrypt.encryptLongText(txt)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(len(contentEncryptList))
	c := RsaEncrypt.decryptLongText(contentEncryptList)
	fmt.Println("decrypt: ", c)
}

func readFile(path, filename string) (string, error) {
	file := path + filename
	d, err := syscall.Open(file, syscall.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer syscall.Close(d)
	var buf [2048]byte
	n, err := syscall.Read(d, buf[:])
	if err != nil {
		return "", err
	}
	contentBuf := make([]byte, n)
	copy(contentBuf, buf[:n])
	return string(contentBuf), nil
}
