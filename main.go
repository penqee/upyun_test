package main

import (
	"bytes"
	"context"
	"filedir/upyun"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/uuid"
)

func main() {
	h := server.Default()

	h.POST("/ping", func(ctx context.Context, c *app.RequestContext) {
		err := os.MkdirAll("tmp/", os.ModePerm)
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

		osFile, err := os.Create("tmp/" + uuid.New().String() + file.Filename)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = io.Copy(osFile, bytes.NewReader(data))
		if err != nil {
			log.Println(err)
			return
		}

		upyun.SaveFile(data, "tmp/", "test/", osFile.Name())

	})

	h.Spin()
}
