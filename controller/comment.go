package controller

import (
	"encoding/json"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment model.Comment `json:"comment,omitempty"`
}

// CommentAction action_type = 1 add comment; action_type = 2 delete comment
func CommentAction(c *gin.Context) {

	// get user information
	UserStr, _ := c.Get("UserStr")
	log.Println("UserStr: ", UserStr)

	var userInfoVar model.User
	if err := json.Unmarshal([]byte(UserStr.(string)), &userInfoVar); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "error: session unmarshal error"})
		return
	}

	// TODO verify
	videoID := c.Query("video_id")
	actionType := c.Query("action_type")
	videoIDNum, _ := strconv.ParseInt(videoID, 10, 64)
	// add comment
	if actionType == "1" {
		commentText := c.Query("comment_text")
		commentVar := &model.Comment{
			UserID:     userInfoVar.ID,
			VideoID:    videoIDNum,
			Content:    commentText,
			User:       userInfoVar,
			CreateData: time.Now().Format("01-02"),
		}
		// db store
		global.DY_DB.Model(&model.Comment{}).Create(&commentVar)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0},
			Comment:  *commentVar,
		})
		return
	}
	// delete comment
	if actionType == "2" {
		commentID := c.Query("comment_id")
		res := global.DY_DB.Delete(&model.Comment{}, "comment_id = ? AND video_id = ?", commentID, videoID)
		if res.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "err: didn't get this comment"})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return
	}

}

// CommentList get comments
func CommentList(c *gin.Context) {

	UserStr, _ := c.Get("UserStr")
	var userInfoVar model.User
	if err := json.Unmarshal([]byte(UserStr.(string)), &userInfoVar); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "error: session unmarshal error"})
		return
	}

	videoID := c.Query("video_id")
	var commentList []model.Comment
	if err := global.DY_DB.Model(&model.Comment{}).Where("video_id = ?", videoID).Preload("User").Find(&commentList).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "db select error"})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
}
