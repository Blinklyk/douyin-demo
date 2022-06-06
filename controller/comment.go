package controller

import (
	"encoding/json"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		tx := global.DY_DB.Begin()
		commentText := c.Query("comment_text")
		// get full user data
		s := service.UserService{}
		returnUser, err := s.GetUserInfo(userInfoVar.ID, userInfoVar.ID)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "error: return user from getUserInfo"})
			return
		}
		commentVar := &model.Comment{
			UserID:     userInfoVar.ID,
			VideoID:    videoIDNum,
			Content:    commentText,
			User:       *returnUser,
			CreateData: time.Now().Format("01-02"),
		}
		// 1. create data in comment table
		if res := tx.Model(&model.Comment{}).Create(&commentVar); res.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "error: create comment in comment table"})
		}

		// corresponding video comment_count + 1
		// 2. update the comment_count column (lock) in video table
		if res := tx.Model(&model.Video{}).Where("id = ?", videoID).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)); res.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "error: update add comment_count in video table"})
			return
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "error: commit update transaction"})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0},
			Comment:  *commentVar,
		})

		return
	}
	// delete comment
	if actionType == "2" {
		tx := global.DY_DB.Begin()
		commentID := c.Query("comment_id")
		// 1. delete comment in comment table
		if res := tx.Delete(&model.Comment{}, "id = ? AND video_id = ?", commentID, videoID); res.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "err: didn't get this comment"})
			return
		}
		// 2. update comment_count - 1 in video table
		if res := tx.Model(&model.Video{}).Where("id = ?", videoID).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", -1)); res.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "error: update subtract comment_count in video table"})
			return
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "error: commit delete transaction"})
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
