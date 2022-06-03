package controller

import (
	"encoding/json"
	"errors"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"strconv"
)

// FavoriteAction directly update db
func FavoriteAction(c *gin.Context) {
	UserStr, _ := c.Get("UserStr")

	log.Println("UserStr: ", UserStr)

	var userInfoVar model.User
	if err := json.Unmarshal([]byte(UserStr.(string)), &userInfoVar); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, CheckUserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "error: session unmarshal error"},
		})
		return
	}
	userID := c.Query("user_id")
	log.Println("userID: ", userID)
	videoID := c.Query("video_id")
	log.Println("videoID: ", videoID)
	// converse to int64 format
	// TODO verify userID nad videoID
	userIDNum, _ := strconv.ParseInt(userID, 10, 64)
	videoIDNum, _ := strconv.ParseInt(videoID, 10, 64)
	//if exist != false {
	//	c.JSON(http.StatusBadRequest, CheckUserInfoResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "error: didn't get video_id in request"},
	//	})
	//	return
	//}
	// action_type determines operationCount

	favoriteInfo := &model.Favorite{
		UserID:  userInfoVar.ID,
		VideoID: videoIDNum,
	}

	actionType := c.Query("action_type")
	actionTypeNum, _ := strconv.ParseInt(actionType, 10, 64)

	if actionTypeNum == 1 {
		// check if add already
		res := global.DY_DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&model.Favorite{})
		if res.RowsAffected != 0 {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "already add favorite this video"})
			return
		}

		// add favorite transaction:
		// 1. create data in user_favorite_video table
		// 2. update favorite_count in videos table
		AddFavorite := func(x *gorm.DB) error {
			tx := global.DY_DB.Begin()

			if err := tx.Create(&favoriteInfo).Error; err != nil {
				log.Println("error when insert u_f_v :", err)
				c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "error: when insert favorite"})
				tx.Rollback()
				return err
			}

			// update the favorite_count column (lock)
			if res := tx.Model(&model.Video{}).Where("id = ?", videoID).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)); res.RowsAffected == 0 {
				tx.Rollback()
				return errors.New("res RowsAffected in videos is 0")
			}

			err := tx.Commit().Error
			return err
		}

		err := AddFavorite(global.DY_DB)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "error when adding favorite: " + err.Error()})
			return
		}
	}

	if actionTypeNum == 2 {
		// check if delete already
		res := global.DY_DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&model.Favorite{})
		if res.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "err: No add favorite this video before"})
			return
		}

		// cancel favorite transaction:
		// 1. delete data in user_favorite_video table  (soft delete)
		// 2. update favorite_count in videos table
		CancelFavorite := func(x *gorm.DB) error {
			tx := global.DY_DB.Begin()

			if err := tx.Delete(&model.Favorite{}, "user_id = ? AND video_id = ?", userIDNum, videoIDNum).Error; err != nil {
				log.Println("error when delete u_f_v :", err)
				tx.Rollback()
				return err
			}

			if res := tx.Model(&model.Video{}).Where("id = ?", videoID).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", -1)); res.RowsAffected == 0 {
				tx.Rollback()
				return errors.New("RowsAffected in video table is 0")
			}

			err := tx.Commit().Error
			return err
		}

		err := CancelFavorite(global.DY_DB)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "error when canceling favorite: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FavoriteList get from favorite table
func FavoriteList(c *gin.Context) {

	// verify
	UserStr, _ := c.Get("UserStr")

	log.Println("UserStr: ", UserStr)

	var userInfoVar model.User
	if err := json.Unmarshal([]byte(UserStr.(string)), &userInfoVar); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "error: session unmarshal error"})
		return
	}

	// find favorite videos from db
	//var userFavoriteVideos model.UserFavoriteVideos
	var favoriteVideoList []model.Video
	var videosID []int64
	// get video_id from conn table first
	res0 := global.DY_DB.Table("dy_favorite").Select("video_id").Where("user_id = ?", userInfoVar.ID).Find(&videosID)
	log.Println("res0.error: ", res0.Error)
	log.Printf("%v\n", videosID)
	// get video details from video table by selecting video_id
	res := global.DY_DB.Model(&model.Video{}).Where("ID in ?", videosID).Find(&favoriteVideoList)
	log.Println("get res RowEffect", res.RowsAffected)

	c.JSON(http.StatusOK, VideoListResponse1{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: favoriteVideoList,
	})
}
