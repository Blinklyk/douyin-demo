package main

import (
	"errors"
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

	// init Viper
	initialize.InitializeConfig()
	// init Redis
	global.DY_REDIS = initialize.InitializeRedis()
	//zap.ReplaceGlobals(global.DY_LOG) // init zap log
	global.DY_DB = initialize.Gorm() // init gorm and connect db
	if global.DY_DB == nil {
		return errors.New("gorm initialize failed")
	} else {
		log.Println("gorm initialize success.")
		initialize.RegisterTables(global.DY_DB) // init tables
	}

	return nil

}
