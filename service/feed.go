package service

import (
	"encoding/json"
	"errors"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/model/request"
	"github.com/RaymondCode/simple-demo/utils"
)

type FeedService struct{}

// FeedWithToken get video information with token
func (fs *FeedService) FeedWithToken(r *request.FeedRequest) (*[]model.Video, error) {
	// parse token
	UserStr, err := utils.RedisParseToken(r.Token)
	if err != nil {
		return nil, errors.New("error: feed parse token")
	}

	var userInfoVar model.User
	if err := json.Unmarshal([]byte(UserStr), &userInfoVar); err != nil {
		return nil, errors.New("error: session unmarshal error")
	}

	// get videos list and traverse it
	videos, err := selectVideos()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(videos); i++ {
		// determine is_favorite value
		var tmp model.Favorite
		if res := global.DY_DB.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userInfoVar.ID, videos[i].ID).First(&tmp); res.RowsAffected != 0 {
			videos[i].IsFavorite = true
		}
		// determine is_follow value
		var temp model.Follow
		if res := global.DY_DB.Model(&model.Follow{}).Where("user_id = ? AND follow_id = ?", userInfoVar.ID, videos[i].UserID).First(&temp); res.RowsAffected != 0 {
			videos[i].User.IsFollow = true
		}
	}

	return &videos, nil

}

// FeedWithoutToken get video information without token
func (fs *FeedService) FeedWithoutToken() (*[]model.Video, error) {
	videos, err := selectVideos()
	if err != nil {
		return nil, err
	}
	return &videos, nil
}

// get videos from db with the latest time
func selectVideos() ([]model.Video, error) {
	// TODO add the latest time
	var videos []model.Video
	if err := global.DY_DB.Preload("User").Order("ID desc").Find(&videos).Error; err != nil {
		return nil, errors.New("error: select without token")
	}
	return videos, nil
}
