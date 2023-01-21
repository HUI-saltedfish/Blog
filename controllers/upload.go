package controllers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Uploader interface {
	upload(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

func UploadGet(c *gin.Context) {
	var (
		err      error
		res      = gin.H{}
		url      string
		uploader Uploader
		file     multipart.File
		fh       *multipart.FileHeader
	)
	defer writeJSON(c, res)
	file, fh, err = c.Request.FormFile("file")
	if err != nil {
		res["message"] = err.Error()
		return
	}

	uploader = &ALiYunUploader{}
	url, err = uploader.upload(file, fh)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["url"] = url
	res["succeed"] = true
}
