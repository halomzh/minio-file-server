package client

import (
	"os"
	"testing"
)

func TestUploadFile(t *testing.T) {
	f, _ := os.Open("/Users/shoufeng/go-work/minio-file/config/check.yml")
	UploadFile("ceshi", "11221", f)
}

func TestFindFileDownloadUrl(t *testing.T) {
	println(FindFileDownloadUrl("ceshi", "11221", "check.yml"))
}

func TestFindFileList(t *testing.T) {
	FindFileList("ceshi", "11221")
}

func TestDeleteFile(t *testing.T) {
	DeleteFile("ceshi", "11221", "check.yml")
}
