package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"minio-file-server/auth"
	"minio-file-server/client"
	"minio-file-server/common"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var defaultLogFile *os.File
var errorLogFile *os.File

type Config struct {
	Port     string   `yaml:"Port"`
	FileType []string `yaml:"FileType"`
}

var ApiConfig = &Config{
	Port:     "9002",
	FileType: []string{},
}

func main() {
	router := gin.Default()

	router.Use(func(context *gin.Context) {
		if !auth.CheckReq(context) {
			context.JSON(http.StatusOK, common.GenFailResult().SetMessage("操作文件失败: 请先登录"))
			context.Abort()
		}
	})

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//上传
	router.POST("/upload", func(c *gin.Context) {

		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["files"]
		//文件类型
		bucketName := c.PostForm("fileBelongType")
		//文件所属
		location := c.PostForm("fileBelongId")

		for _, file := range files {
			log.Println(file.Filename)
			isAcceptFileType := false
			for _, ft := range ApiConfig.FileType {
				if strings.HasSuffix(file.Filename, ft) {
					isAcceptFileType = true
					break
				}
			}
			if !isAcceptFileType {
				c.JSON(200, common.GenFailResult().SetMessage("无效文件类型: "+file.Filename))
				c.Abort()
				return
			}
			fileNameTemp := os.TempDir() + "/" + file.Filename
			c.SaveUploadedFile(file, fileNameTemp)
			fileTemp, _ := os.Open(fileNameTemp)
			client.UploadFile(bucketName, location, fileTemp)
			os.Remove(fileNameTemp)
		}

		c.JSON(200, common.GenSuccessResult())
	})

	//下载地址
	router.POST("/findDownloadUrl", func(c *gin.Context) {

		//文件类型
		bucketName := c.PostForm("fileBelongType")
		//文件所属
		location := c.PostForm("fileBelongId")
		//文件名称
		fileName := c.PostForm("fileName")
		url := client.FindFileDownloadUrl(bucketName, location, fileName)

		c.JSON(200, common.GenSuccessResult().SetData(url))
	})

	//删除
	router.POST("/delete", func(c *gin.Context) {

		//文件类型
		bucketName := c.PostForm("fileBelongType")
		//文件所属
		location := c.PostForm("fileBelongId")
		//文件名称
		fileName := c.PostForm("fileName")
		flag := client.DeleteFile(bucketName, location, fileName)
		if flag {
			c.JSON(200, common.GenSuccessResult())
		} else {
			c.JSON(200, common.GenFailResult())
		}

	})

	//查询
	router.POST("/findFileList", func(c *gin.Context) {

		//文件类型
		bucketName := c.PostForm("fileBelongType")
		//文件所属
		location := c.PostForm("fileBelongId")

		fileNameList := client.FindFileList(bucketName, location)

		c.JSON(200, common.GenSuccessResult().SetData(fileNameList))
	})

	router.Run(":9002")
}

func init() {
	ApiConfig.init()

	logFileName := "./log/" + time.Now().Format("2006-01-02")
	for i := 0; ; i++ {
		logFileNameTemp := logFileName + "（" + strconv.Itoa(i) + "）"
		_, err := os.Stat(logFileNameTemp + "default.log")
		if common.CheckError(err) {
			defaultLogFile, _ = os.Create(logFileNameTemp + "default.log")
			errorLogFile, _ = os.Create(logFileNameTemp + "error.log")
			break
		}
	}
	gin.DefaultWriter = io.MultiWriter(defaultLogFile)
	gin.DefaultErrorWriter = io.MultiWriter(errorLogFile)

}

func (config *Config) init() {
	configByte, err := ioutil.ReadFile("./config/api.yml")
	common.CheckError(err)
	err = yaml.Unmarshal(configByte, config)
	common.CheckError(err)
}
