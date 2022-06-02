package global

import (
	"github.com/RaymondCode/simple-demo/config"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/sessions"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

var (
	DY_DB            *gorm.DB
	DY_DBList        map[string]*gorm.DB
	DY_CONFIG        config.Server
	DY_LOG           *zap.Logger
	DY_JWTMW         *jwt.GinJWTMiddleware
	DY_REDIS         *redis.Client
	DY_SESSION_STORE sessions.Store
)

const (
	SECRETKEY    = "secrete key"
	DY_OSSDOMAIN = "http://rceumi5re.bkt.gdipper.com/"
)
