package upyun

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/upyun/go-sdk/v3/upyun"
)

var UpYun *upyun.UpYun

func NewUpYun() {
	UpYun = upyun.NewUpYun(
		&upyun.UpYunConfig{
			Bucket:   "w2-domtok",
			Operator: "penqee",
			Password: "operator-password",
		},
	)
}

func uploadFile(src, dest string) error {
	return UpYun.Put(&upyun.PutObjectConfig{
		Path:      dest,
		UseMD5:    true,
		LocalPath: src,
	})
}

// func GetImageUrl(uri string) (string, error) {
// 	etime := strconv.FormatInt(time.Now().Unix()+config.Upyun.TokenTimeout, 10)
// 	sign := utils.MD5(strings.Join([]string{config.Upyun.TokenSecret, etime, uri}, "&"))
// 	url := fmt.Sprintf("%s%s?_upt=%s%s", config.Upyun.UssDomain, utils.UriEncode(uri), sign[12:20], etime)
// 	return url, nil
// }

func SaveFile(data []byte, tmpFile, destDir, filename string) error {

	out, err := os.Create(tmpFile + uuid.New().String() + filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		return err
	}
	go uploadToUpYun(tmpFile, destDir)
	return nil
}

func uploadToUpYun(tempDir string, destDir string) {
	files, err := os.ReadDir(tempDir)
	if err != nil {
		logger.Errorf("upyun ReadDir failed, err: %v", err)
		return
	}

	for _, file := range files {
		filename := file.Name()
		src := filepath.Join(tempDir, filename)
		dest := filepath.Join(destDir, filename)

		err = uploadFile(src, dest)
		if err != nil {
			logger.Errorf("upyun upload file failed, err: %v", err)
		}

		err = os.Remove(src)
		if err != nil {
			logger.Errorf("upyun remove file failed, err: %v", err)
		}
	}
}

func Split(str string) string {
	strs := strings.Split(str, "/")
	return strs[len(strs)-1]
}
