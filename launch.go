package main

import "meetroom/server/server"

func main() {
	server.NewServer()

	//  初始化mysql数据库连接
	server.InitDatabase();
	// 初始化redis数据库连接
	server.InitRedis();

	// register Mid ware 注册中间件
	server.UseIOCMid(server.JWTVerifyIOC)

	//register Application 注册应用

	//启动服务
	server.Build(":8848")
}