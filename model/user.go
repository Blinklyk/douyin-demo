package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	// TODO 改成嵌入结构体
	ID             int64          `gorm:"primarykey"` // 主键ID
	CreatedAt      time.Time      // 创建时间
	UpdatedAt      time.Time      // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                            // 删除时间
	Name           string         `json:"name,omitempty" gorm:"default:testName"`    // TODO
	FollowCount    int64          `json:"follow_count,omitempty" gorm:"default:0"`   // 关注数
	FollowerCount  int64          `json:"follower_count,omitempty" gorm:"default:0"` // 粉丝数
	IsFollow       bool           `json:"is_follow,omitempty" gorm:"default:false"`  // 当前用户是否关注
	Username       string         `json:"username" gorm:"comment:username" `         // 登录账号
	Password       string         `json:"password" gorm:"comment:password"`          // 登录密码
	Videos         []Video        `json:"videos" gorm:"comment:pulished videos"`     // 视频列表
	FavoriteVideos []Video        `gorm:"many2many:favorite"`                        // 点赞视频列表
	Followers      []User         `gorm:"many2many:follow"`                          // 粉丝列表
}
