package dynamic

import (
	"meetroom/server/IOC"
	"meetroom/server/manage"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)


func dynamicView(db *gorm.DB,rC *redis.Conn)manage.Viewer{
	v:=manage.QuickNewViewer("/:divPath/:urlPath",db,rC,
	)

	return v
}

type DynamicAction struct {
	DivPath IOC.Value `ioc:"from:path;to:string;name:divPath"`
	URIPath IOC.Value `ioc:"from:path;to:string;name:urlPath"`
}
