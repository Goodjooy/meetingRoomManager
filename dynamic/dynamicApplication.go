package dynamic

import "meetroom/server/manage"

//TODO:包内共享的配置文件信息  

func dynamicApplication() manage.Application{
	app:=manage.NewApplication("/","dynamic","")

	return app
}	