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

		filess, err := c.MultipartForm()
		if err != nil {
			log.Println(err)
			return
		}

		for index, file := range filess.File["file"] {
			src, err := file.Open()
			if err != nil {
				log.Println(err)
				return
			}

			defer src.Close()
			log.Println(index)
			data, err := ioutil.ReadAll(src)
			if err != nil {
				log.Println(err)
				return
			}

			upyun.SaveFile(data, tempFile, destDir, file.Filename)
		}

		// if  files, ok := filess.File["file"]; ok {

		// 	// log.Println(index)
		// 	// file, err := c.FormFile("file")
		// 	// if err != nil {
		// 	// 	log.Println(err)
		// 	// 	return
		// 	// }

		// }

	})

	h.Spin()
}
