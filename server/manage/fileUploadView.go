package manage

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadFileViewHandle(c *gin.Context, targetKey, appName string)string {
	file, err := c.FormFile(targetKey)
	if err != nil {
		c.String(http.StatusBadRequest,"上传模式异常")
	}

	fileName := file.Filename
	v := strings.Split(fileName, ".")
	fileEx := v[len(v)-1]
	fileName = UUIDGenerate()

	saveFileName := fmt.Sprintf("%s.%s", fileName, fileEx)
	savePath ,_:= filepath.Abs(fmt.Sprintf("%s/%s", mediaSaveRoot, appName))
	if _,err:= os.Stat(savePath); err!=nil{
		os.MkdirAll(savePath,os.ModePerm)
	}
	saveFilePath:=fmt.Sprintf("%s/%s",savePath,saveFileName)
	
	err=c.SaveUploadedFile(file, saveFilePath)
	if err != nil {
		c.String(500,"上传出错")
	}

	return fmt.Sprintf("%s/%s/%s",MediaURLRoot,appName,saveFileName)
}
