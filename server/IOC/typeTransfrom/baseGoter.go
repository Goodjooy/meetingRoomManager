package typeTransfrom

import "github.com/gin-gonic/gin"

func fromPostForm(c *gin.Context) func(key string) (data interface{}, exist bool) {
	return func(key string) (data interface{}, exist bool) {
		v := c.PostForm(key)
		exist = true
		if v == "" {
			exist = false
		}
		return v, exist
	}
}
func fromURLPath(c *gin.Context) func(key string) (data interface{}, exist bool) {
	return func(key string) (data interface{}, exist bool) {
		v := c.Param(key)
		exist = true
		if v == "" {
			exist = false
		}
		return v, exist
	}
}
func fromQuery(c *gin.Context) func(key string) (data interface{}, exist bool) {
	return func(key string) (data interface{}, exist bool) {
		v := c.Query(key)
		exist = true
		if v == "" {
			exist = false
		}
		return v, exist
	}
}
func fromContext(c *gin.Context) func(key string) (data interface{}, exist bool) {
	return c.Get
}
func fromCookie(c *gin.Context) func(key string) (data interface{}, exist bool) {
	return func(key string) (data interface{}, exist bool) {
		v, err := c.Cookie(key)
		exist = err == nil
		return v, exist
	}
}
func fromMulit(c *gin.Context) func(key string) (data interface{}, exist bool) {
	return func(key string) (data interface{}, exist bool) {
		
		v, err := c.FormFile(key)
		exist = err == nil
		return v, exist
	}
}

