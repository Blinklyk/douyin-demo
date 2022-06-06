package request

import "mime/multipart"

type PublishRequest struct {
	data  *multipart.FileHeader `json:"data"`
	Token string                `json:"token" form:"token"`
	Title string                `json:"title" form:"title"`
}
