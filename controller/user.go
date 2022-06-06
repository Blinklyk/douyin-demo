package controller

import (
	"encoding/json"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/model/request"
	"github.com/RaymondCode/simple-demo/model/response"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type CheckUserInfoResponse struct {
	Response
	UserInfo model.User `json:"user"`
}

func Register(c *gin.Context) {

	// 把校验接受数据以及校验放在结构体Register上
	var r request.RegisterRequest
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser := &model.User{Username: r.Username, Password: r.Password, FollowCount: 0, FollowerCount: 0}
	var userService = service.UserService{}
	err, userReturn := userService.Register(*newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RegisterResponse{
			Response: response.Response{
				StatusCode: -1,
				StatusMsg:  "failed: create register user " + err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, response.RegisterResponse{
			Response: response.Response{
				StatusCode: 0,
				StatusMsg:  "success: create register user" + err.Error(),
			},
			UserId: userReturn.ID,
			Token:  userReturn.Username,
		})
	}

}

func Login(c *gin.Context) {

	type LoginVar struct {
		Username string `json:"username" gorm:"not null; comment:username for register;" form:"username"`
		Password string `json:"password" gorm:"not null; comment:password for register" form:"password"`
	}
	var l LoginVar
	if err := c.ShouldBind(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{Username: l.Username, Password: l.Password}
	var loginService = service.UserService{}
	userReturn, tokenStr, err := loginService.Login(*user)
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error tokenStr is empty": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response.LoginResponse{
			Response: response.Response{
				StatusCode: -1,
				StatusMsg:  "failed: login in" + err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, response.LoginResponse{
			Response: response.Response{
				StatusCode: 0,
				StatusMsg:  "success: login in",
			},
			UserId: userReturn.ID,
			Token:  tokenStr,
		})
	}

	return
}

// UserInfo get userInfo from db
func UserInfo(c *gin.Context) {

	//jwt version
	UserStr, _ := c.Get("UserStr")

	log.Println("UserStr: ", UserStr)

	var userInfoVar model.User
	if err := json.Unmarshal([]byte(UserStr.(string)), &userInfoVar); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, CheckUserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "error: session unmarshal error"},
		})
		return
	}

	// TODO check
	if len(userInfoVar.Name) < 3 {
		c.JSON(http.StatusOK, CheckUserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "userName len less then 3"},
		})
		return
	}
	// get user info from db
	var checkUserInfoService = service.UserService{}
	returnUser, err := checkUserInfoService.GetUserInfo(userInfoVar.ID, userInfoVar.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "error: db select"})
		return
	}

	c.JSON(http.StatusOK, CheckUserInfoResponse{
		Response: Response{StatusCode: 0},
		UserInfo: *returnUser,
	})
	return

	////session version
	//// get user from redis
	//userID := c.Query("user_id")
	//session := sessions.Default(c)
	//jsonUser := session.Get(userID)
	//log.Println("jsonUser : ", jsonUser)
	//log.Println(c.GetHeader("Cookie"))
	//log.Println(c.GetHeader("Host"))
	//log.Println(c.GetHeader("Connection"))
	//
	////userInfoVar := &userInfoVar{}
	//userInfoVar := &model.User{}
	//
	//err := json.Unmarshal(jsonUser.([]byte), userInfoVar)
	//if err != nil {
	//	c.JSON(http.StatusOK, CheckUserInfoResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "Unmarshal from session failed"},
	//	})
	//	return
	//}
	//
	//if len(userInfoVar.Name) < 0 {
	//	c.JSON(http.StatusOK, CheckUserInfoResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "userName len is 0"},
	//	})
	//	return
	//}
	//
	//
	//c.JSON(http.StatusOK, CheckUserInfoResponse{
	//	Response: Response{StatusCode: 0},
	//	UserInfo: *userInfoVar,
	//})
	//return

}
