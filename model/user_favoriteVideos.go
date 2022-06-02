package model

import (
	"github.com/RaymondCode/simple-demo/global"
	"gorm.io/gorm"
	"log"
	"time"
)

type UserFavoriteVideos struct {
	UserID    int64          `gorm:"primaryKey"`
	VideoID   int64          `gorm:"primaryKey"`
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

func setRelation() {
	if err := global.DY_DB.SetupJoinTable(&User{}, "FavoriteVideos", &UserFavoriteVideos{}); err != nil {
		log.Println(err)
		panic("set join table error")
	}

}
