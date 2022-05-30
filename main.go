package main

import (
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/initialize"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {

	if err := Init(); err != nil {
		os.Exit(-1)
	}
	r := gin.Default()

	initRouter(r)
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run(":" + global.App.Config.App.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func Init() error {

	// 初始化Viper
	initialize.InitializeConfig()
	// 初始化Redis
	global.DY_REDIS = initialize.InitializeRedis()
	//zap.ReplaceGlobals(global.DY_LOG) // 初始化zap日志
	global.DY_DB = initialize.Gorm() // gorm连接数据库
	if global.DY_DB == nil {
		log.Println("gorm initialize failed.")
	} else {
		log.Println("gorm initialize success.")
		initialize.RegisterTables(global.DY_DB) // 初始化表
		// 程序结束前关闭数据库链接
		//db, _ := global.DY_DB.DB()
		//defer db.Close()
	}

	//if err := repository.Init(); err != nil {
	//	return err
	//}
	return nil

}
