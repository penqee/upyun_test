package main

import (
	"context"
	"filedir/upyun"
	"io/ioutil"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()

	h.POST("/ping", func(ctx context.Context, c *app.RequestContext) {

		const (
			tempFile = "tmp/"
			destDir  = "test/"
		)

		err := os.MkdirAll(tempFile, os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}
		upyun.NewUpYun()
		file, err := c.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}
		src, err := file.Open()
		if err != nil {
			log.Println(err)
			return
		}

		defer src.Close()

		data, err := ioutil.ReadAll(src)
		if err != nil {
			log.Println(err)
			return
		}

		upyun.SaveFile(data, tempFile, destDir, file.Filename)

	})

	h.Spin()
}
