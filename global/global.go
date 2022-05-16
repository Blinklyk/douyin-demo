package glabal

import (
	"gorm.io/gorm"
)

var (
	DY_DB     *gorm.DB
	DY_DBList map[string]*gorm.DB
)