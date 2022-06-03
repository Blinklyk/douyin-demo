package model

type follow struct {
	ID         int64 `gorn:"primarykey"`
	FollowerID int64
}
