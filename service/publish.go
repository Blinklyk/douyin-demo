package service

import (
	"errors"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/model/request"
	"github.com/RaymondCode/simple-demo/utils"
	"log"
	"time"
)

type PublishService struct{}

// PublishAction publish the video to oss and get the url
func (ps *PublishService) PublishAction(u *model.User, r *request.PublishRequest, filePath string) error {
	title := r.Title

	// upload the file to oss
	ret := utils.UploadFile(filePath)
	// get the url from oss
	VideoUrl := global.DY_OSS_DOMAIN + ret.Key
	log.Println("Publish video url: " + VideoUrl)

	publishVideo := &model.Video{
		UserID:        u.ID,
		PlayUrl:       VideoUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		PublishTime:   time.Now(),
		Title:         title,
		IsFavorite:    false,
	}

	if result := global.DY_DB.Model(&model.Video{}).Create(&publishVideo); result.RowsAffected == 0 {
		return errors.New("publish error")
	}
	return nil
}

// PublishList return the publishing video list
func (ps *PublishService) PublishList(r *request.PublishListRequest) (publishVideos []model.Video, err error) {
	if err := global.DY_DB.Where("user_id = ?", r.UserID).Preload("User").Order("ID desc").Find(&publishVideos).Error; err != nil {
		return nil, err
	}

	return

}
