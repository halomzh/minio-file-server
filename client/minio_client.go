package client

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"minio-file-server/common"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	Endpoint        string `yaml:"Endpoint"`
	AccessKeyID     string `yaml:"AccessKeyID"`
	SecretAccessKey string `yaml:"SecretAccessKey"`
}

var MinioConfig *Config = &Config{
	Endpoint:        "",
	AccessKeyID:     "",
	SecretAccessKey: "",
}

var MinioClient *minio.Client

func init() {
	MinioConfig.init()

	// 初使化 minio client对象。
	minioClient, err := minio.New(MinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioConfig.AccessKeyID, MinioConfig.SecretAccessKey, ""),
		Secure: false,
	})
	common.CheckError(err)
	MinioClient = minioClient
	log.Printf("%#v\n", minioClient) // minioClient初使化成功
}

func (config *Config) init() {
	configByte, err := ioutil.ReadFile("./config/minio.yml")
	common.CheckError(err)
	err = yaml.Unmarshal(configByte, config)
	common.CheckError(err)
}

func UploadFile(bucketName, location string, file *os.File) {
	//检查bucket是否存在
	isBucketExists, err := MinioClient.BucketExists(context.Background(), bucketName)
	common.CheckError(err)
	if !isBucketExists {
		err = MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "", ObjectLocking: false})
		common.CheckError(err)
	}

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	uploadInfo, err := MinioClient.PutObject(context.Background(), bucketName, location+"/"+fileStat.Name(), file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	common.CheckError(err)
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)

}

func FindFileDownloadUrl(bucketName, location, fileName string) string {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fileName+"\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := MinioClient.PresignedGetObject(context.Background(), bucketName, location+"/"+fileName, time.Second*24*60*60, reqParams)
	common.CheckError(err)
	fmt.Println("Successfully generated presigned URL", presignedURL)

	return presignedURL.String()
}

func FindFileList(bucketName, location string) []string {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    location,
		Recursive: true,
	})

	fileNameList := make([]string, 0)

	index := 0
	for ch := range objectCh {
		nameSpTemp := strings.Split(ch.Key, "/")
		fileNameList = append(fileNameList, nameSpTemp[len(nameSpTemp)-1])
		index++
	}
	for object := range objectCh {
		if object.Err != nil {
			common.CheckError(object.Err)
		}
		fmt.Println(object)
	}

	return fileNameList
}

func DeleteFile(bucketName, location, fileName string) bool {
	opts := minio.RemoveObjectOptions{}
	err := MinioClient.RemoveObject(context.Background(), bucketName, location+"/"+fileName, opts)

	return !common.CheckError(err)
}
