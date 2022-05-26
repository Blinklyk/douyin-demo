package controller

import (
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

var videos []model.Video

// Feed same demo video list for every request
func Feed(c *gin.Context) {

	result := global.DY_DB.Preload("User").Order("ID desc").Find(&videos)
	if result.RowsAffected == 0 {
		log.Println("0 videos query from database")
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
