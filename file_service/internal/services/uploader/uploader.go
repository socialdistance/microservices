package uploader

import (
	"log/slog"
	"mime/multipart"
)

const (
	ods = "ods"
	csv = "csv"
)

type Uploader struct {
	log        *slog.Logger
	TypeFolder string
	Size       int
	Extension  string
}

func NewUploader(log *slog.Logger, typeFolder, extension string, size int) *Uploader {

	return &Uploader{
		log:        log,
		TypeFolder: typeFolder,
		Size:       size,
		Extension:  extension,
	}
}

func (u *Uploader) SaveFile(filename, typeFolder string, mediaFile multipart.File) error {

	return nil
}
