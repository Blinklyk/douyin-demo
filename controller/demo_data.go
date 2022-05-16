package controller

var DemoVideos = []Video{
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 99,
		CommentCount:  3,
		IsFavorite:    false,
	},
	{
		Id:            2,
		Author:        DemoUser,
		PlayUrl:       "http://vfx.mtime.cn/Video/2019/02/04/mp4/190204084208765161.mp4",
		CoverUrl:      "https://upload.wikimedia.org/wikipedia/zh/2/2b/Captain_Marvel_Poster.jpg",
		FavoriteCount: 100,
		CommentCount:  2,
		IsFavorite:    false,
	},
}

var DemoComments = []Comment{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
