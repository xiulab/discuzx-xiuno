package lfile

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
)

// MD5 文件 MD5
func MD5(file *multipart.FileHeader) (md5Str string, err error) {
	f, err := file.Open()
	if err != nil {
		return
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
