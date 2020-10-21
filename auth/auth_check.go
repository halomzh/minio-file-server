package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minio-file-server/common"
	"net/http"
)

var AuthConfig *Config = &Config{Url: "", HeaderName: ""}

type Config struct {
	IsOpen     bool   `yaml:"IsOpen"`
	Url        string `yaml:"Url"`
	HeaderName string `yaml:"HeaderName"`
}

func CheckReq(c *gin.Context) bool {
	if !AuthConfig.IsOpen {
		return true
	}
	accessToken := c.GetHeader(AuthConfig.HeaderName)
	req, err := http.NewRequest(http.MethodGet, AuthConfig.Url, nil)
	common.CheckError(err)
	req.Header.Add(AuthConfig.HeaderName, accessToken)

	resp, _ := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	r := &common.Result{}
	err = json.Unmarshal(body, r)
	common.CheckError(err)
	if r.Code != 500 {
		return true
	}

	return false
}

func (config *Config) init() {
	configByte, err := ioutil.ReadFile("./config/check.yml")
	common.CheckError(err)
	err = yaml.Unmarshal(configByte, config)
	common.CheckError(err)
}

func init() {
	AuthConfig.init()
}
