package models

import "mime/multipart"

type Media struct {
	MediaFile multipart.FileHeader `json:"media_file,omitempty" validate:"required"`
}
