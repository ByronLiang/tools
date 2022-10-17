package wf

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var ErrFileNameIsDir = errors.New("filename is dir")

func fileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if info.IsDir() {
		return false, ErrFileNameIsDir
	}
	return true, nil
}

func CheckAndCreateFilePath(filename string) error {
	exist, isDirErr := fileExists(filename)
	if isDirErr != nil {
		return isDirErr
	}
	if !exist {
		return createDirAndFile(filepath.Dir(filename), filename)
	}
	return nil
}

func createDirAndFile(dir string, filename string) error {
	// 检测目录是否存在
	_, StatErr := os.Stat(dir)
	if os.IsNotExist(StatErr) {
		mkdirErr := os.Mkdir(dir, 0744)
		if mkdirErr != nil {
			return mkdirErr
		}
	}
	// 目录存在, 只需创建文件
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return f.Close()
}

func GenWatchPath(dir string, filename string) (string, error) {
	fpath := filepath.Join(dir, filename)
	err := CheckAndCreateFilePath(fpath)
	if err != nil {
		return "", err
	}
	return fpath, err
}

func recreateFile(file string) error {
	if err := os.Remove(file); err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	return f.Close()
}

func BuildTempFile(prefix string) (string, error) {
	// temp file
	fObj, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}
	tempFileName := fObj.Name()
	fObj.Close()
	return tempFileName, nil
}

func BuildTempDirAndFile(dir string, filename string, total int) (fileDir string, filenameList []string, err error) {
	fileDir, err = ioutil.TempDir("", dir)
	if err != nil {
		return
	}
	filenameList = make([]string, 0, total)
	var fileObj *os.File
	for i := 0; i < total; i++ {
		fileObj, err = ioutil.TempFile(fileDir, filename)
		if err != nil {
			continue
		}
		filenameList = append(filenameList, fileObj.Name())
		fileObj.Close()
	}
	return
}
