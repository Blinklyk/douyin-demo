package controller

import (
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
type VideoListResponse1 struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//token := c.PostForm("token")
	//
	//if _, exist := usersLoginInfo[token]; !exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}
	//
	//data, err := c.FormFile("data")
	//if err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//
	//filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]
	//finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	//saveFile := filepath.Join("./public/", finalName)
	//if err := c.SaveUploadedFile(data, saveFile); err != nil {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, Response{
	//	StatusCode: 0,
	//	StatusMsg:  finalName + " uploaded successfully",
	//})

	userID, exist := c.Get("ID")
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Didn't get the ID from token"})
		return
	}
	title := c.PostForm("title")
	// save the file at local host
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	filename := filepath.Base(data.Filename)
	saveFile := filepath.Join("public/", filename)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// upload the file to oss
	ret := utils.UploadFile("public/" + filename)
	// get the url from oss
	VideoUrl := global.DY_OSSDOMAIN + ret.Key
	log.Println(VideoUrl)

	publishVideo := &model.Video{
		UserID:        userID.(int64),
		PlayUrl:       VideoUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		PublishTime:   time.Now(),
		Title:         title,
	}

	result := global.DY_DB.Model(&model.Video{}).Create(&publishVideo)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "发布成功！",
		})
		return
	}
}

// PublishList

var publishVideos []model.Video

func PublishList(c *gin.Context) {

	userID, exist := c.Get("ID")
	if !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Didn't get the ID from token"})
		return
	}

	result := global.DY_DB.Where("user_id = ?", userID).Preload("User").Order("ID desc").Find(&publishVideos)
	if result.RowsAffected == 0 {
		log.Println("0 videos query from database")
	}

	c.JSON(http.StatusOK, VideoListResponse1{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: publishVideos,
	})
}
