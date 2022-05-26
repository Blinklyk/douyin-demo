package service

import (
	"errors"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

type UserService struct{}

func (us *UserService) Register(user model.User) (err error, newUser model.User) {
	// 校验 查询数据库中是否有此用户(高级查询)
	var u model.User
	if !errors.Is(global.DY_DB.Model(&model.User{}).Where("username = ?", user.Username).First(&u).Error, gorm.ErrRecordNotFound) {
		return errors.New("this username is registered already"), user
	}
	// 雪花算法生成新的id
	var node, _ = utils.NewWorker(1)
	newID := node.GetId()
	user.ID = newID
	// 密码加密
	user.Password = utils.BcryptHash(user.Password)
	// 添加到数据库
	err = global.DY_DB.Create(&user).Error
	return err, user
}

func (us *UserService) Login(user model.User) (returnUser *model.User, tokenStr string, err error) {
	// TODO 校验
	// 查询 账号密码是否正确
	var u model.User
	if errors.Is(global.DY_DB.Model(&model.User{}).Where("username = ?", user.Username).First(&u).Error, gorm.ErrRecordNotFound) {
		return nil, "", errors.New("user doesn't exsit")
	}
	if ok := utils.BcryptCheck(user.Password, u.Password); !ok {
		return nil, "", errors.New("password error")
	}
	// 生成token
	tokenStr, err = utils.GenToken(u.ID)
	if err != nil {
		return nil, "", err
	}
	return &u, tokenStr, nil

}

func (us *UserService) CheckUserInfo(id int64) (returnUser *model.User, err error) {
	var u model.User
	if errors.Is(global.DY_DB.Model(&model.User{}).Where("id = ?", id).First(&u).Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user doesn't exsit")
	}

	return &u, nil

}
