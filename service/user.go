package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/RaymondCode/simple-demo/global"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
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

	// jwt version
	// TODO 校验
	// 查询 账号密码是否正确
	var u model.User

	// get user form db
	if errors.Is(global.DY_DB.Model(&model.User{}).Where("username = ?", user.Username).First(&u).Error, gorm.ErrRecordNotFound) {
		return nil, "", errors.New("user doesn't exist")
	}
	if ok := utils.BcryptCheck(user.Password, u.Password); !ok {
		return nil, "", errors.New("password error")
	}
	log.Printf("get User data from db : %v", u)
	// gen token
	tokenStr, err = utils.GenToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	// store user data into redis
	jsonU, err := json.Marshal(u)
	if err != nil {
		return nil, "", errors.New("json marshal error")
	}
	// redis key: "login:session:"+tokenStr, value: user TTL: 30min
	res := global.DY_REDIS.Set(context.Background(), "login:session:"+tokenStr, jsonU, time.Minute*30)
	log.Println("res.String() user set to redis:", res)
	return &u, tokenStr, nil

	//// session + redis version
	//// TODO check format
	//var u model.User
	//if errors.Is(global.DY_DB.Model(&model.User{}).Where("username = ?", user.Username).First(&u).Error, gorm.ErrRecordNotFound) {
	//	return nil, "", errors.New("user doesn't exist")
	//}
	//if ok := utils.BcryptCheck(user.Password, u.Password); !ok {
	//	return nil, "", errors.New("password error")
	//}
	//// 生成session key : userID, value: user
	//jsonU, err := json.Marshal(u)
	//if err != nil {
	//	return nil, "", errors.New("json marshal error")
	//}
	//session.Set(u.ID, jsonU)
	//err = session.Save()

	if err != nil {
		return nil, "", err
	}
	return &u, strconv.FormatInt(u.ID, 10), nil

}

func (us *UserService) CheckUserInfo(id int64) (returnUser *model.User, err error) {
	var u model.User
	if errors.Is(global.DY_DB.Model(&model.User{}).Where("id = ?", id).First(&u).Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user doesn't exist")
	}

	return &u, nil

}
