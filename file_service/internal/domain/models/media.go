package models

import "mime/multipart"

type Media struct {
	MediaFile multipart.File `json:"media_file,omitempty" validate:"required"`
}

type URL struct {
	Url string `json:"url,omitempty" validate:"required"`
}
