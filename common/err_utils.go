package common

import "github.com/gin-gonic/gin"

func CheckError(err error) bool {
	if err != nil {
		_, _ = gin.DefaultErrorWriter.Write([]byte(err.Error() + "\n"))
		return true
	}
	return false
}
