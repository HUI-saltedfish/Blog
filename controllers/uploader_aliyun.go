package controllers

import (
	"Blog/system"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cihub/seelog"
)

type ALiYunUploader struct{}

func (a *ALiYunUploader) upload(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	endpoint := system.GetConfiguration().ALiYunEndpoint
	accessKeyId := system.GetConfiguration().ALiYunAccessKey
	accessKeySecret := system.GetConfiguration().ALiYunSecretkey
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		seelog.Errorf("oss.New err: %v", err)
		return "", err
	}

	// 指定自己的bucket
	bucketName := system.GetConfiguration().QiniuBucketName
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		seelog.Errorf("client.Bucket err: %v", err)
		return "", err
	}

	src, err := fileHeader.Open()
	if err != nil {
		seelog.Errorf("file.Open err: %v", err)
		return "", err
	}
	defer src.Close()

	// 上传文件流 并定义上传的文件夹
	folderName := time.Now().Format("2006-01-02")
	fileTmpPath := filepath.Join("uploads", folderName) + "/" + fileHeader.Filename
	err = bucket.PutObject(fileTmpPath, src)
	if err != nil {
		seelog.Errorf("bucket.PutObject err: %v", err)
		return "", err
	}

	// 获取文件的url
	url := "https://" + bucketName + "." + endpoint + "/" + fileTmpPath
	return url, nil
}
