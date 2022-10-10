package wf

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

var ErrFilePathIsDir = errors.New("filepath is dir")

func fileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if info.IsDir() {
		return false, ErrFilePathIsDir
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
		log.Println("remove-file err", err.Error())
	}
	f, err := os.Create(file)
	if err != nil {
		log.Println("recreateFile err")
		return err
	}
	return f.Close()
}
