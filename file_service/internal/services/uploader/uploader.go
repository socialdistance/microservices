package uploader

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
)

type Uploader struct {
	log      *slog.Logger
	filePath string
}

func NewUploader(log *slog.Logger, filePath string) *Uploader {

	return &Uploader{
		log:      log,
		filePath: filePath,
	}
}

func (u *Uploader) SaveFile(folderType string, mediaFile multipart.FileHeader) error {
	const op = "uploader.SaveFile"

	writerRecovery := bytes.NewBufferString(u.filePath + "recovery/" + mediaFile.Filename)
	writerCreate := bytes.NewBufferString(u.filePath + "create/" + mediaFile.Filename)

	if _, err := os.Stat(u.filePath); os.IsNotExist(err) {
		return fmt.Errorf("Directory for save files does not exist %s: %s", op, err)
	} else {
		fmt.Println("Directory exists")
	}

	src, err := mediaFile.Open()
	if err != nil {
		return fmt.Errorf("error open file %s: %s", op, err.Error())
	}
	defer src.Close()

	if folderType == "recovery" {
		err := selectFolderForSaveFile(writerRecovery.String(), src)
		if err != nil {
			return fmt.Errorf("error open file %s: %s", op, err)
		}
	}

	if folderType == "create" {
		err := selectFolderForSaveFile(writerCreate.String(), src)
		if err != nil {
			return fmt.Errorf("error open file %s: %s", op, err)
		}
	}

	return nil
}

func selectFolderForSaveFile(writer string, src multipart.File) error {
	const op = "uploader.selectFolderForSaveFile"

	file, err := os.OpenFile(writer, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error open file %s: %s", op, err.Error())
	}
	defer file.Close()

	if _, err = io.Copy(file, src); err != nil {
		return fmt.Errorf("can't save file to destination create directory")
	}

	return nil
}
