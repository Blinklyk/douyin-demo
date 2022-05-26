package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type CheckUserInfoResponse struct {
	Response
	UserInfo model.User `json:"user"`
}

func Register(c *gin.Context) {

	//username := c.Query("username")
	//password := c.Query("password")
	//token := username + password
	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else {
	//	atomic.AddInt64(&userIdSequence, 1)
	//	newUser := User{
	//		Id:   userIdSequence,
	//		Name: username,
	//	}
	//	usersLoginInfo[token] = newUser
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   userIdSequence,
	//		Token:    username + password,
	//	})
	//}
	// 把校验接受数据以及校验放在结构体Register上
	var r model.Register
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser := &model.User{ID: 1, Username: r.Username, Password: r.Password}
	var userService = service.UserService{}
	err, userReturn := userService.Register(*newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "failed: create register user " + err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success: create register user ",
			},
			UserId: userReturn.ID,
			Token:  userReturn.Username + userReturn.Password,
		})
	}

}

func Login(c *gin.Context) {
	//username := c.Query("username")
	//password := c.Query("password")
	//
	//token := username + password
	//
	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   user.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}

	var l model.Login
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
		c.JSON(http.StatusBadRequest, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "failed: login in" + err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success: login in",
			},
			UserId: userReturn.ID,
			Token:  tokenStr,
		})
	}

	return
}

// UserInfo jwt中放userID 从数据库查找返回user数据
func UserInfo(c *gin.Context) {

	//if user, exist := usersLoginInfo[tokenStr]; exist {
	//	c.JSON(http.StatusOK, CheckUserInfoResponse{
	//		Response: Response{StatusCode: 0},
	//		User:     user,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, CheckUserInfoResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist1"},
	//	})
	//}

	userID, exist := c.Get("ID")
	log.Println(userID)
	log.Println(exist)

	var userInfoVar struct {
		ID            int64  `json:"id"`
		Name          string `json:"name,omitempty"`
		FollowCount   int64  `json:"follow_count,omitempty"`
		FollowerCount int64  `json:"follower_count,omitempty"`
		IsFollow      bool   `json:"is_follow,omitempty"`
	}
	if err := c.BindQuery(&userInfoVar); err != nil {
		c.JSON(http.StatusOK, CheckUserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "error: bind error"},
		})
		return
	}

	if len(userInfoVar.Name) < 0 {
		c.JSON(http.StatusOK, CheckUserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "userName len is 0"},
		})
		return
	}
	var checkUserInfoService = service.UserService{}
	returnUser, err := checkUserInfoService.CheckUserInfo(userID.(int64))
	if err != nil {
		c.JSON(http.StatusOK, CheckUserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist1"},
		})
		return
	} else {
		c.JSON(http.StatusOK, CheckUserInfoResponse{
			Response: Response{StatusCode: 0},
			UserInfo: *returnUser,
		})
		return
	}

}
