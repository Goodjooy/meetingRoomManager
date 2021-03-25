package manage

import "github.com/gin-gonic/gin"

func AddFileManage(server *gin.Engine){
	files:=server.Group(MediaURLRoot)
	files.GET("/*file",func(c *gin.Context) {
		file:=c.Param("file")
		c.File(mediaSaveRoot+ file)
	})
}

