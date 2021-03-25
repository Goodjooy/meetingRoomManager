package manage

import (
	"crypto/sha256"
	"fmt"
	"meetroom/server/IOC"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type administer interface {
	GetTargetModel(appName, modelName string) interface{}
}
type quickAsign interface {
	GetURLPattern() string
	GetHandles() map[string][]gin.HandlerFunc
}

/*
Application 后端的每个单独应用
*/
type Application struct {
	URLPattern string
	name       string

	viewers           []Viewer
	viewerURLPatterns []string

	useAddon []gin.HandlerFunc

	models []interface{}
}

func NewApplication(URLPattern, appName, elseData string) Application {
	app := Application{URLPattern: URLPattern, name: appName}

	return app
}

func methodLimitaion(f func(c *gin.Context), supportMethods []string) gin.HandlerFunc {
	temp := func(c *gin.Context) {
		isSupport := false
		method := c.Request.Method
		for _, v := range supportMethods {
			if v == method {
				isSupport = true
				break
			}
		}

		if isSupport {
			f(c)
		} else {
			c.String(http.StatusNotFound, fmt.Sprintf("the method \"%s\" not support", method))
		}
	}

	return temp
}

func (app Application) GetAllModels() []interface{} {
	return app.models
}
func (app Application) GetAppName() string {
	return app.name
}
func (app *Application) AsignAddon(addons ...gin.HandlerFunc) {
	app.useAddon = append(app.useAddon, addons...)
}
func (app *Application) AsignAddonIOC(db *gorm.DB, rC *redis.Conn, addons...interface{}) {
	var fns []gin.HandlerFunc
	for _, v := range addons {
		t := IOC.ToIOC(v)
		f := IOC.DoIOC(t, db,rC)
		fns = append(fns, f)
	}

	app.useAddon = append(app.useAddon, fns...)
}
func (app *Application) AsignModels(models ...interface{}) {
	app.models = append(app.models, models...)
}

func (app *Application) AsignViewer(viewers ...Viewer) {
	for _, viewer := range viewers {
		viewerURLPattern := viewer.URLPattern
		exist := false

		for _, pattern := range app.viewerURLPatterns {
			if pattern == viewerURLPattern {
				exist = true
			}
		}

		if !exist {
			app.viewers = append(app.viewers, viewer)
			app.viewerURLPatterns = append(app.viewerURLPatterns, viewerURLPattern)
		}
	}
}

func (app Application) AsignApplication(server *gin.Engine, db *gorm.DB) Application {
	group := server.Group(app.URLPattern)

	group.Use(app.useAddon...)
	viewers := app.viewers
	for _, v := range viewers {
		handles := v.handle
		for medhod, fn := range handles {
			group.Handle(string(medhod), v.URLPattern, fn...)
		}
	}

	for _, model := range app.models {
		if !db.HasTable(model) {
			db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(model)
		}
	}
	return app
}

const numStart = '0'
const lowAlphaStart = 'a'
const upAlphaStart = 'A'

func DateSHA256Hash(passwd string) string {
	
	hash := sha256.Sum256([]byte(passwd))

	var hashRe []byte

	for _, v := range hash {
		t := v % (10 + 26 + 26)
		if t < 10 {
			t = t + numStart
		} else if t < 26+10 {
			t = t + lowAlphaStart - 10
		} else {
			t = t + upAlphaStart - 10 - 26
		}
		hashRe = append(hashRe, t)
	}

	return string(hashRe)
}
func UUIDGenerate() string {
	value := uuid.NewV4()

	return value.String()
}
